package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type DivisionFactory struct {
	service *services.DivisionService
	handler *handlers.DivisionHandler
	routes  *routes.DivisionRoutes
}

func NewDivisionFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *DivisionFactory {
	service := services.NewDivisionService(
		cont.DivisionRepo,
		cont.TenantRepo,
		cont.Logger,
	)
	handler := handlers.NewDivisionHandler(
		mw.ContextMiddleware,
		service,
		cont.Logger,
	)
	return &DivisionFactory{
		service: service,
		handler: handler,
		routes: routes.NewDivisionRoutes(
			handler,
			mw.ContextMiddleware,
			mw.AuthMiddleware,
			mw.TenantMiddleware,
		),
	}
}

func (f *DivisionFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
