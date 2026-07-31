package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type SuperadminFactory struct {
	service *services.SuperadminService
	handler *handlers.SuperadminHandler
	routes  *routes.SuperadminRoutes
}

func NewSuperadminFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *SuperadminFactory {
	service := services.NewSuperadminService(cont.TenantRepo, cont.Logger)
	handler := handlers.NewSuperadminHandler(mw.ContextMiddleware, service, cont.Logger)
	return &SuperadminFactory{
		service: service,
		handler: handler,
		routes:  routes.NewSuperadminRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware),
	}
}

func (f *SuperadminFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
