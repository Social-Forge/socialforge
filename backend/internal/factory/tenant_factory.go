package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type TenantFactory struct {
	service *services.TenantService
	handler *handlers.TenantHandler
	routes  *routes.TenantRoutes
}

func NewTenantFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *TenantFactory {
	service := services.NewTenantService(
		cont.TenantRepo,
		cont.Logger,
		cont.MinioClient,
	)
	handler := handlers.NewTenantHandler(
		mw.ContextMiddleware,
		service,
		cont.Logger,
	)
	return &TenantFactory{
		service: service,
		handler: handler,
		routes: routes.NewTenantRoutes(
			handler,
			mw.ContextMiddleware,
			mw.AuthMiddleware,
			mw.TenantMiddleware,
			mw.CSRFMiddleware,
		),
	}
}
func (f *TenantFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
