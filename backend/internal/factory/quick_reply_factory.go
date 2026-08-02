package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type QuickReplyFactory struct {
	service *services.QuickReplyService
	handler *handlers.QuickReplyHandler
	routes  *routes.QuickReplyRoutes
}

func NewQuickReplyFactory(cont *dependencies.Container, mw *MiddlewareFactory) *QuickReplyFactory {
	service := services.NewQuickReplyService(cont.QuickReplyRepo, cont.TenantRepo, cont.Logger)
	handler := handlers.NewQuickReplyHandler(mw.ContextMiddleware, service, cont.Logger)
	return &QuickReplyFactory{
		service: service,
		handler: handler,
		routes:  routes.NewQuickReplyRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *QuickReplyFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
