package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type ContactFactory struct {
	service *services.ContactService
	handler *handlers.ContactHandler
	routes  *routes.ContactRoutes
}

func NewContactFactory(cont *dependencies.Container, mw *MiddlewareFactory) *ContactFactory {
	service := services.NewContactService(cont.ContactRepo, cont.Logger)
	handler := handlers.NewContactHandler(mw.ContextMiddleware, service, cont.Logger)
	return &ContactFactory{
		service: service,
		handler: handler,
		routes:  routes.NewContactRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *ContactFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
