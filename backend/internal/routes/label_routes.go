package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type LabelRoutes struct {
	path      string
	handler   *handlers.LabelHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewLabelRoutes(
	handler *handlers.LabelHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *LabelRoutes {
	return &LabelRoutes{path: "/labels", handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *LabelRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Post("/", r.handler.Create)
	protected.Put("/:id", r.handler.Update)
	protected.Delete("/:id", r.handler.Delete)

	// Conversation <-> label attachment
	protected.Get("/conversation/:conversationId", r.handler.ListForConversation)
	protected.Post("/conversation/:conversationId/attach", r.handler.Attach)
	protected.Delete("/conversation/:conversationId/label/:labelId", r.handler.Detach)
}
