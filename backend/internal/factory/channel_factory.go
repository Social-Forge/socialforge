package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type ChannelFactory struct {
	service *services.ChannelService
	handler *handlers.ChannelHandler
	routes  *routes.ChannelRoutes
}

func NewChannelFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *ChannelFactory {
	service := services.NewChannelService(
		cont.ChannelRepo,
		cont.DivisionRepo,
		cont.TenantRepo,
		cont.AutoResponseRepo,
		services.BuildConnectors(cont.Config),
		cont.Config.App.URL,
		"",
		cont.Logger,
	)
	handler := handlers.NewChannelHandler(mw.ContextMiddleware, service, cont.Logger)
	return &ChannelFactory{
		service: service,
		handler: handler,
		routes: routes.NewChannelRoutes(
			handler,
			mw.ContextMiddleware,
			mw.AuthMiddleware,
			mw.TenantMiddleware,
		),
	}
}

func (f *ChannelFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
