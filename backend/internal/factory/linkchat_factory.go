package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type LinkchatFactory struct {
	service *services.LinkchatService
	handler *handlers.LinkchatHandler
	routes  *routes.LinkchatRoutes
}

func NewLinkchatFactory(cont *dependencies.Container, mw *MiddlewareFactory) *LinkchatFactory {
	service := services.NewLinkchatService(cont.DivisionRepo, cont.Logger)
	handler := handlers.NewLinkchatHandler(mw.ContextMiddleware, service, cont.Logger)
	return &LinkchatFactory{
		service: service,
		handler: handler,
		routes:  routes.NewLinkchatRoutes(handler, mw.ContextMiddleware),
	}
}

func (f *LinkchatFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
