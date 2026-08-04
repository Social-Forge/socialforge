package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type AIAgentRoutes struct {
	path      string
	handler   *handlers.AIAgentHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewAIAgentRoutes(
	handler *handlers.AIAgentHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *AIAgentRoutes {
	return &AIAgentRoutes{path: "/ai-agents", handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *AIAgentRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Post("/", r.handler.Create)
	protected.Get("/:id", r.handler.Get)
	protected.Put("/:id", r.handler.Update)
	protected.Delete("/:id", r.handler.Delete)
}
