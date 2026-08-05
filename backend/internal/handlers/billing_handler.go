package handlers

import (
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BillingHandler struct {
	ctxinject *middlewares.ContextMiddleware
	billing   *services.BillingService
	subs      *services.SubscriptionService
	logger    *zap.Logger
}

func NewBillingHandler(ctxinject *middlewares.ContextMiddleware, billing *services.BillingService, subs *services.SubscriptionService, logger *zap.Logger) *BillingHandler {
	return &BillingHandler{ctxinject: ctxinject, billing: billing, subs: subs, logger: logger}
}

func (h *BillingHandler) tenantUUID(c fiber.Ctx) (uuid.UUID, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(tid)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func (h *BillingHandler) isOwner(c fiber.Ctx) bool {
	names, ok := c.Locals("role_name").([]string)
	if !ok {
		return false
	}
	for _, n := range names {
		if n == entity.RoleTenantOwner {
			return true
		}
	}
	return false
}

// Checkout starts a payment (owner only).
func (h *BillingHandler) Checkout(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantUUID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can manage billing", nil)
	}
	var req dto.CheckoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	payerEmail, _ := c.Locals("email").(string)
	result, err := h.billing.Checkout(ctx, tenantID, payerEmail, &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Checkout created successfully", result)
}

// GetSubscription returns the tenant's current subscription + plan + limits.
func (h *BillingHandler) GetSubscription(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantUUID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	view, err := h.subs.GetCurrent(ctx, tenantID)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Subscription retrieved successfully", view)
}

// ListInvoices returns the tenant's invoice history.
func (h *BillingHandler) ListInvoices(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantUUID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	invoices, err := h.billing.ListInvoices(ctx, tenantID)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Invoices retrieved successfully", invoices)
}

// Webhook is the public provider callback endpoint. Authenticity is verified
// inside the gateway (token/signature); always returns 200 on a handled event so
// the provider does not retry indefinitely, 4xx only on auth failure.
func (h *BillingHandler) Webhook(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	provider := c.Params("provider")
	body := c.Body()
	headers := map[string]string{
		"X-Callback-Token":         c.Get("X-Callback-Token"),
		"Paypal-Transmission-Id":   c.Get("Paypal-Transmission-Id"),
		"Paypal-Transmission-Time": c.Get("Paypal-Transmission-Time"),
		"Paypal-Transmission-Sig":  c.Get("Paypal-Transmission-Sig"),
		"Paypal-Cert-Url":          c.Get("Paypal-Cert-Url"),
		"Paypal-Auth-Algo":         c.Get("Paypal-Auth-Algo"),
	}
	if err := h.billing.HandleWebhook(ctx, provider, headers, body); err != nil {
		h.logger.Warn("billing webhook failed", zap.String("provider", provider), zap.Error(err))
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "ok", nil)
}
