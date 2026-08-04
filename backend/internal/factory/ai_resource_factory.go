package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type AIResourceFactory struct {
	service *services.AIResourceService
	handler *handlers.AIResourceHandler
	routes  *routes.AIResourceRoutes
}

func NewAIResourceFactory(cont *dependencies.Container, mw *MiddlewareFactory) *AIResourceFactory {
	service := services.NewAIResourceService(
		cont.AIAgentRepo,
		cont.AIKnowledgeRepo,
		cont.AIPlaybookRepo,
		cont.AIAssetRepo,
		cont.Logger,
	)
	handler := handlers.NewAIResourceHandler(mw.ContextMiddleware, service, cont.Logger)
	return &AIResourceFactory{
		service: service,
		handler: handler,
		routes:  routes.NewAIResourceRoutes(handler, mw.ContextMiddleware, mw.AuthMiddleware, mw.TenantMiddleware),
	}
}

func (f *AIResourceFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
