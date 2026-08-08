package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SubscriptionService reconciles the billing subscription record with the
// tenant's enforced entitlements (tenant.Max* + ai_credits + subscription_plan).
type SubscriptionService struct {
	subRepo    repository.SubscriptionRepository
	planRepo   repository.PlanRepository
	tenantRepo repository.TenantRepository
	addonRepo  repository.SubscriptionAddonRepository
	logger     *zap.Logger
}

func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
	tenantRepo repository.TenantRepository,
	addonRepo repository.SubscriptionAddonRepository,
	logger *zap.Logger,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:    subRepo,
		planRepo:   planRepo,
		tenantRepo: tenantRepo,
		addonRepo:  addonRepo,
		logger:     logger,
	}
}

// CurrentView bundles a tenant's subscription with its resolved plan.
type CurrentView struct {
	Subscription *entity.Subscription `json:"subscription"`
	Plan         *entity.Plan         `json:"plan"`
	Tenant       *entity.Tenant       `json:"tenant"`
}

// GetCurrent returns the tenant's active subscription + plan (+ tenant limits).
func (s *SubscriptionService) GetCurrent(ctx context.Context, tenantID uuid.UUID) (*CurrentView, error) {
	view := &CurrentView{}
	var err error
	view.Tenant, err = s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	tctx := repository.WithTenantID(ctx, tenantID)
	_ = s.subRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		view.Subscription, _ = s.subRepo.FindActiveByTenant(txCtx, tenantID)
		return nil
	})
	if view.Subscription != nil {
		view.Plan, _ = s.planRepo.FindByID(ctx, view.Subscription.PlanID)
	}
	if view.Plan == nil {
		// Fall back to the tenant's plan code (e.g. free, never checked out).
		view.Plan, _ = s.planRepo.FindByCode(ctx, view.Tenant.SubscriptionPlan)
	}
	return view, nil
}

// ActivatePlan upserts the subscription to active for `months` and applies the
// plan's entitlements to the tenant. Called on successful payment.
func (s *SubscriptionService) ActivatePlan(ctx context.Context, tenantID uuid.UUID, plan *entity.Plan, months int) error {
	if months <= 0 {
		months = 1
	}
	now := time.Now()
	end := now.AddDate(0, months, 0)

	tctx := repository.WithTenantID(ctx, tenantID)
	err := s.subRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.subRepo.FindActiveByTenant(txCtx, tenantID)
		if err != nil {
			return err
		}
		if existing == nil {
			return s.subRepo.Create(txCtx, &entity.Subscription{
				TenantID:           tenantID,
				PlanID:             plan.ID,
				Status:             entity.SubscriptionStatusActive,
				CurrentPeriodStart: entity.NewNullTime(now),
				CurrentPeriodEnd:   entity.NewNullTime(end),
			})
		}
		existing.PlanID = plan.ID
		existing.Status = entity.SubscriptionStatusActive
		existing.CurrentPeriodStart = entity.NewNullTime(now)
		existing.CurrentPeriodEnd = entity.NewNullTime(end)
		existing.CancelAtPeriodEnd = false
		_, err = s.subRepo.Update(txCtx, existing)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %w", err)
	}

	// Apply entitlements to the tenant (the enforcement source of truth).
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return err
	}
	applyPlanToTenant(tenant, plan)
	if _, err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return fmt.Errorf("failed to apply plan to tenant: %w", err)
	}
	s.logger.Info("subscription activated",
		zap.String("tenant_id", tenantID.String()),
		zap.String("plan", plan.Code),
		zap.Time("period_end", end))
	return nil
}

// GrantAddon records an addon purchase and applies its immediate effect (AI
// credits are added to the tenant balance; slot addons raise effective quota,
// resolved at enforcement time in 6E).
func (s *SubscriptionService) GrantAddon(ctx context.Context, tenantID uuid.UUID, addonType string, quantity int, meta entity.JSONMap) error {
	tctx := repository.WithTenantID(ctx, tenantID)
	if err := s.addonRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.addonRepo.Create(txCtx, &entity.SubscriptionAddon{
			TenantID: tenantID,
			Type:     addonType,
			Quantity: quantity,
			Meta:     meta,
		})
	}); err != nil {
		return fmt.Errorf("failed to record addon: %w", err)
	}

	if addonType == entity.AddonTypeAICredits && quantity > 0 {
		tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
		if err != nil {
			return err
		}
		tenant.AiCredits += quantity
		if _, err := s.tenantRepo.Update(ctx, tenant); err != nil {
			return fmt.Errorf("failed to add ai credits: %w", err)
		}
	}
	s.logger.Info("addon granted", zap.String("tenant_id", tenantID.String()),
		zap.String("type", addonType), zap.Int("quantity", quantity))
	return nil
}

// SweepExpired downgrades subscriptions whose period has ended: the subscription
// row is marked expired and the tenant is reverted to the free plan's limits.
// Runs cross-tenant (worker), so it operates pool-level without RLS ctx.
func (s *SubscriptionService) SweepExpired(ctx context.Context) (int, error) {
	now := time.Now()
	expired, err := s.subRepo.ListExpired(ctx, now)
	if err != nil {
		return 0, err
	}
	if len(expired) == 0 {
		return 0, nil
	}
	freePlan, err := s.planRepo.FindByCode(ctx, entity.PlanFree)
	if err != nil {
		return 0, fmt.Errorf("sweep: free plan missing: %w", err)
	}

	downgraded := 0
	for _, sub := range expired {
		// Skip if a renewal already pushed the period forward.
		if sub.CurrentPeriodEnd.Valid && sub.CurrentPeriodEnd.Time.After(now) {
			continue
		}
		tctx := repository.WithTenantID(ctx, sub.TenantID)
		sub.Status = entity.SubscriptionStatusExpired
		if err := s.subRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
			_, uerr := s.subRepo.Update(txCtx, sub)
			return uerr
		}); err != nil {
			s.logger.Warn("sweep: failed to expire subscription", zap.String("id", sub.ID.String()), zap.Error(err))
			continue
		}
		tenant, err := s.tenantRepo.FindByID(ctx, sub.TenantID)
		if err != nil {
			continue
		}
		applyPlanToTenant(tenant, freePlan)
		if _, err := s.tenantRepo.Update(ctx, tenant); err != nil {
			s.logger.Warn("sweep: failed to downgrade tenant", zap.String("tenant_id", sub.TenantID.String()), zap.Error(err))
			continue
		}
		downgraded++
		s.logger.Info("subscription expired -> downgraded to free", zap.String("tenant_id", sub.TenantID.String()))
	}
	return downgraded, nil
}

// applyPlanToTenant maps a plan's feature entitlements onto the tenant's enforced
// limits. Missing feature keys leave the current value unchanged.
func applyPlanToTenant(t *entity.Tenant, plan *entity.Plan) {
	f := plan.Features
	set := func(field *int, key string) {
		if v, ok := f.Int(key); ok {
			*field = v
		}
	}
	set(&t.MaxDivisions, "divisions")
	set(&t.MaxAgents, "agents")
	set(&t.MaxQuickReplies, "quick_replies")
	set(&t.MaxWahaWhatsApp, "waha_whatsapp")
	set(&t.MaxMetaWhatsApp, "meta_whatsapp")
	set(&t.MaxMetaMessenger, "meta_messenger")
	set(&t.MaxInstagram, "instagram")
	set(&t.MaxTelegram, "telegram")
	set(&t.MaxWebChat, "webchat")
	set(&t.MaxLinkChat, "linkchat")
	if credits, ok := f.Int("ai_credits"); ok {
		t.AiCredits = credits // monthly allotment (reset on activation/renewal)
	}
	t.SubscriptionPlan = plan.Code
	t.SubscriptionStatus = entity.StatusActive
	// An active plan supersedes any trial window — clear it so the tenant is not
	// gated by a stale trial_ends_at.
	t.TrialEndsAt = entity.NullTime{}
}
