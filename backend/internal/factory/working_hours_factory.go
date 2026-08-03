package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type WorkingHoursFactory struct {
	service *services.WorkingHoursService
	handler *handlers.WorkingHoursHandler
	routes  *routes.WorkingHoursRoutes
}

func NewWorkingHoursFactory(cont *dependencies.Container, mw *MiddlewareFactory) *WorkingHoursFactory {
	service := services.NewWorkingHoursService(cont.WorkingHoursRepo, cont.Logger)
	handler := handlers.NewWorkingHoursHandler(mw.ContextMiddleware, service, cont.Logger)
	return &WorkingHoursFactory{
		service: service,
		handler: handler,
		routes:  routes.NewWorkingHoursRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *WorkingHoursFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
