package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/payments"
	"github/socialforge/internal/infra/repository"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Addon unit prices (whole IDR). ai_credits is priced per credit.
const (
	addonPriceChannelSlot  = 50000
	addonPriceAgentSlot    = 50000
	addonPricePerAICredit  = 10 // Rp10 per credit (e.g. 10.000 credits = Rp100.000)
)

// BillingService orchestrates checkout (invoice + gateway) and webhook handling
// (mark paid + apply entitlements), provider-agnostic via the Gateway interface.
type BillingService struct {
	invoiceRepo repository.InvoiceRepository
	eventRepo   repository.PaymentEventRepository
	planRepo    repository.PlanRepository
	subSvc      *SubscriptionService
	gateways    map[string]payments.Gateway
	baseURL     string
	logger      *zap.Logger
}

func NewBillingService(
	invoiceRepo repository.InvoiceRepository,
	eventRepo repository.PaymentEventRepository,
	planRepo repository.PlanRepository,
	subSvc *SubscriptionService,
	gateways map[string]payments.Gateway,
	baseURL string,
	logger *zap.Logger,
) *BillingService {
	return &BillingService{
		invoiceRepo: invoiceRepo,
		eventRepo:   eventRepo,
		planRepo:    planRepo,
		subSvc:      subSvc,
		gateways:    gateways,
		baseURL:     strings.TrimRight(baseURL, "/"),
		logger:      logger,
	}
}

// CheckoutResult is returned to the client to redirect the payer.
type CheckoutResult struct {
	Invoice     *entity.Invoice `json:"invoice"`
	CheckoutURL string          `json:"checkout_url"`
}

// Checkout creates a pending invoice and a hosted-checkout session on the chosen
// provider, returning the URL to redirect the payer to.
func (s *BillingService) Checkout(ctx context.Context, tenantID uuid.UUID, payerEmail string, req *dto.CheckoutRequest) (*CheckoutResult, error) {
	gateway := s.gateways[req.Provider]
	if gateway == nil {
		return nil, fmt.Errorf("payment provider %q not available", req.Provider)
	}

	amount, description, purpose, err := s.resolveCharge(ctx, req)
	if err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, fmt.Errorf("nothing to charge (amount is zero)")
	}

	inv := &entity.Invoice{
		TenantID:    tenantID,
		Status:      entity.InvoiceStatusPending,
		Amount:      amount,
		Currency:    "IDR",
		Description: description,
		Purpose:     purpose,
		Provider:    req.Provider,
		ExpiresAt:   entity.NewNullTime(time.Now().Add(24 * time.Hour)),
	}
	tctx := repository.WithTenantID(ctx, tenantID)
	if err := s.invoiceRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.invoiceRepo.Create(txCtx, inv)
	}); err != nil {
		return nil, err
	}

	created, err := gateway.CreateInvoice(ctx, payments.CreateInvoiceParams{
		ExternalID:  inv.ID.String(),
		Number:      inv.Number,
		Amount:      amount,
		Currency:    "IDR",
		Description: description,
		PayerEmail:  payerEmail,
		SuccessURL:  s.baseURL + "/billing/success?invoice=" + inv.ID.String(),
		FailureURL:  s.baseURL + "/billing/failed?invoice=" + inv.ID.String(),
	})
	if err != nil {
		// Mark the invoice failed so it isn't left dangling as pending.
		_ = s.invoiceRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
			return s.invoiceRepo.SetStatus(txCtx, inv.ID, entity.InvoiceStatusFailed)
		})
		return nil, fmt.Errorf("gateway checkout failed: %w", err)
	}

	_ = s.invoiceRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.invoiceRepo.SetProviderInfo(txCtx, inv.ID, created.ProviderInvoiceID, created.CheckoutURL)
	})
	inv.ProviderInvoiceID = entity.NewNullString(created.ProviderInvoiceID)
	inv.CheckoutURL = entity.NewNullString(created.CheckoutURL)

	s.logger.Info("checkout created",
		zap.String("tenant_id", tenantID.String()),
		zap.String("provider", req.Provider),
		zap.Int("amount", amount),
		zap.Int("invoice_number", inv.Number))
	return &CheckoutResult{Invoice: inv, CheckoutURL: created.CheckoutURL}, nil
}

// resolveCharge computes the amount, description and purpose for a checkout.
func (s *BillingService) resolveCharge(ctx context.Context, req *dto.CheckoutRequest) (int, string, entity.JSONMap, error) {
	switch req.Kind {
	case entity.InvoicePurposeSubscription:
		if req.PlanCode == "" {
			return 0, "", nil, fmt.Errorf("plan_code is required for a subscription")
		}
		plan, err := s.planRepo.FindByCode(ctx, req.PlanCode)
		if err != nil {
			return 0, "", nil, err
		}
		months := req.Months
		if months <= 0 {
			months = 1
		}
		amount := plan.Price * months
		desc := fmt.Sprintf("Langganan %s (%d bulan)", plan.Name, months)
		purpose := entity.JSONMap{"kind": entity.InvoicePurposeSubscription, "plan_code": plan.Code, "months": months}
		return amount, desc, purpose, nil

	case entity.InvoicePurposeAddon:
		if req.Quantity <= 0 {
			return 0, "", nil, fmt.Errorf("quantity is required for an addon")
		}
		var unit int
		switch req.AddonType {
		case entity.AddonTypeChannelSlot:
			unit = addonPriceChannelSlot
		case entity.AddonTypeAgentSlot:
			unit = addonPriceAgentSlot
		case entity.AddonTypeAICredits:
			unit = addonPricePerAICredit
		default:
			return 0, "", nil, fmt.Errorf("unknown addon_type %q", req.AddonType)
		}
		amount := unit * req.Quantity
		desc := fmt.Sprintf("Add-on %s x%d", req.AddonType, req.Quantity)
		purpose := entity.JSONMap{"kind": entity.InvoicePurposeAddon, "addon_type": req.AddonType, "quantity": req.Quantity}
		return amount, desc, purpose, nil
	}
	return 0, "", nil, fmt.Errorf("unknown checkout kind %q", req.Kind)
}

// HandleWebhook authenticates a provider callback, records it, and — on a
// pending→paid transition — applies the invoice purpose (activate plan / grant
// addon). Idempotent: a redelivered "paid" callback is a no-op.
func (s *BillingService) HandleWebhook(ctx context.Context, provider string, headers map[string]string, body []byte) error {
	gateway := s.gateways[provider]
	if gateway == nil {
		return fmt.Errorf("payment provider %q not available", provider)
	}
	res, err := gateway.VerifyWebhook(headers, body)
	if err != nil {
		return err
	}

	// Resolve the invoice: prefer our external id (round-tripped), else provider id.
	var inv *entity.Invoice
	if res.ExternalID != "" {
		if id, perr := uuid.Parse(res.ExternalID); perr == nil {
			inv, err = s.invoiceRepo.FindByID(ctx, id)
		}
	}
	if inv == nil && res.ProviderInvoiceID != "" {
		inv, err = s.invoiceRepo.FindByProviderInvoiceID(ctx, provider, res.ProviderInvoiceID)
	}
	if err != nil || inv == nil {
		return fmt.Errorf("webhook invoice not found (provider=%s ext=%s pid=%s)", provider, res.ExternalID, res.ProviderInvoiceID)
	}

	tctx := repository.WithTenantID(ctx, inv.TenantID)

	// Audit the event.
	_ = s.eventRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.eventRepo.Create(txCtx, &entity.PaymentEvent{
			TenantID:   inv.TenantID,
			InvoiceID:  &inv.ID,
			Provider:   provider,
			EventType:  res.EventType,
			ExternalID: entity.NewNullString(res.EventID),
			Payload:    entity.JSONMap{"raw": string(body)},
		})
	})

	if res.Expired {
		_ = s.invoiceRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
			return s.invoiceRepo.SetStatus(txCtx, inv.ID, entity.InvoiceStatusExpired)
		})
		return nil
	}
	if !res.Paid {
		return nil // non-terminal event; nothing to apply
	}

	// Pending -> paid transition (idempotent).
	var transitioned bool
	if err := s.invoiceRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		transitioned, err = s.invoiceRepo.MarkPaid(txCtx, inv.ID, time.Now())
		return err
	}); err != nil {
		return err
	}
	if !transitioned {
		s.logger.Info("webhook: invoice already paid, skipping", zap.Int("invoice_number", inv.Number))
		return nil
	}

	return s.applyPurpose(ctx, inv)
}

// applyPurpose fulfils a paid invoice: activate a plan or grant an addon.
func (s *BillingService) applyPurpose(ctx context.Context, inv *entity.Invoice) error {
	kind, _ := inv.Purpose["kind"].(string)
	switch kind {
	case entity.InvoicePurposeSubscription:
		code, _ := inv.Purpose["plan_code"].(string)
		months, _ := inv.Purpose.Int("months")
		plan, err := s.planRepo.FindByCode(ctx, code)
		if err != nil {
			return fmt.Errorf("apply subscription: %w", err)
		}
		return s.subSvc.ActivatePlan(ctx, inv.TenantID, plan, months)
	case entity.InvoicePurposeAddon:
		addonType, _ := inv.Purpose["addon_type"].(string)
		qty, _ := inv.Purpose.Int("quantity")
		meta := entity.JSONMap{}
		if addonType == entity.AddonTypeAICredits {
			meta["credits"] = qty
		}
		return s.subSvc.GrantAddon(ctx, inv.TenantID, addonType, qty, meta)
	}
	return fmt.Errorf("unknown invoice purpose kind %q", kind)
}

// ListInvoices returns a tenant's invoice history.
func (s *BillingService) ListInvoices(ctx context.Context, tenantID uuid.UUID) ([]*entity.Invoice, error) {
	tctx := repository.WithTenantID(ctx, tenantID)
	var out []*entity.Invoice
	err := s.invoiceRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		var err error
		out, err = s.invoiceRepo.ListByTenant(txCtx, tenantID)
		return err
	})
	return out, err
}
