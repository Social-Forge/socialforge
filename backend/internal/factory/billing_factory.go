package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/infra/payments"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type BillingFactory struct {
	handler *handlers.BillingHandler
	routes  *routes.BillingRoutes
}

func NewBillingFactory(cont *dependencies.Container, mw *MiddlewareFactory) *BillingFactory {
	subsService := services.NewSubscriptionService(
		cont.SubscriptionRepo, cont.PlanRepo, cont.TenantRepo, cont.AddonRepo, cont.Logger,
	)
	gateways := payments.BuildGateways(&cont.Config.Payment)
	billingService := services.NewBillingService(
		cont.InvoiceRepo, cont.PaymentEventRepo, cont.PlanRepo, subsService,
		gateways, cont.Config.App.ClientOrigin, cont.Logger,
	)
	handler := handlers.NewBillingHandler(mw.ContextMiddleware, billingService, subsService, cont.Logger)
	return &BillingFactory{
		handler: handler,
		routes:  routes.NewBillingRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *BillingFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
