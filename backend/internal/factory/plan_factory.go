package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type PlanFactory struct {
	service *services.PlanService
	handler *handlers.PlanHandler
	routes  *routes.PlanRoutes
}

func NewPlanFactory(cont *dependencies.Container, mw *MiddlewareFactory) *PlanFactory {
	service := services.NewPlanService(cont.PlanRepo, cont.Logger)
	handler := handlers.NewPlanHandler(mw.ContextMiddleware, service, cont.Logger)
	return &PlanFactory{
		service: service,
		handler: handler,
		routes:  routes.NewPlanRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware),
	}
}

func (f *PlanFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
