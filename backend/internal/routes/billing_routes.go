package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type BillingRoutes struct {
	path      string
	handler   *handlers.BillingHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewBillingRoutes(
	handler *handlers.BillingHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *BillingRoutes {
	return &BillingRoutes{path: "/billing", handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *BillingRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	// Public provider callbacks (authenticity verified inside the gateway).
	route.Post("/webhooks/:provider", r.handler.Webhook)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())
	protected.Post("/checkout", r.handler.Checkout)
	protected.Get("/subscription", r.handler.GetSubscription)
	protected.Get("/invoices", r.handler.ListInvoices)
}
