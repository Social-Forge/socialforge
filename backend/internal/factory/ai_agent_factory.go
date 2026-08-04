package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type AIAgentFactory struct {
	service *services.AIAgentService
	handler *handlers.AIAgentHandler
	routes  *routes.AIAgentRoutes
}

func NewAIAgentFactory(cont *dependencies.Container, mw *MiddlewareFactory) *AIAgentFactory {
	service := services.NewAIAgentService(cont.AIAgentRepo, cont.Logger)
	handler := handlers.NewAIAgentHandler(mw.ContextMiddleware, service, cont.Logger)
	return &AIAgentFactory{
		service: service,
		handler: handler,
		routes:  routes.NewAIAgentRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *AIAgentFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
