package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

// AIResourceRoutes mounts an AI agent's children (knowledge, playbooks, assets)
// nested under the agent: /ai-agents/protected/:id/{knowledge|playbooks|assets}.
type AIResourceRoutes struct {
	handler   *handlers.AIResourceHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewAIResourceRoutes(
	handler *handlers.AIResourceHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *AIResourceRoutes {
	return &AIResourceRoutes{handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *AIResourceRoutes) RegisterRoutes(parent fiber.Router) {
	agent := parent.Group("/ai-agents/protected/:id")
	agent.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	kn := agent.Group("/knowledge")
	kn.Get("/", r.handler.ListKnowledge)
	kn.Post("/", r.handler.CreateKnowledge)
	kn.Put("/:childId", r.handler.UpdateKnowledge)
	kn.Delete("/:childId", r.handler.DeleteKnowledge)

	pb := agent.Group("/playbooks")
	pb.Get("/", r.handler.ListPlaybooks)
	pb.Post("/", r.handler.CreatePlaybook)
	pb.Put("/:childId", r.handler.UpdatePlaybook)
	pb.Delete("/:childId", r.handler.DeletePlaybook)

	as := agent.Group("/assets")
	as.Get("/", r.handler.ListAssets)
	as.Post("/", r.handler.CreateAsset)
	as.Put("/:childId", r.handler.UpdateAsset)
	as.Delete("/:childId", r.handler.DeleteAsset)
}
