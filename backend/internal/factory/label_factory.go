package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type LabelFactory struct {
	service *services.LabelService
	handler *handlers.LabelHandler
	routes  *routes.LabelRoutes
}

func NewLabelFactory(cont *dependencies.Container, mw *MiddlewareFactory) *LabelFactory {
	service := services.NewLabelService(cont.LabelRepo, cont.CentrifugoClient, cont.Logger)
	handler := handlers.NewLabelHandler(mw.ContextMiddleware, service, cont.Logger)
	return &LabelFactory{
		service: service,
		handler: handler,
		routes:  routes.NewLabelRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *LabelFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
