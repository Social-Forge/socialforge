package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type ConversationRoutes struct {
	path      string
	handler   *handlers.ConversationHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewConversationRoutes(
	handler *handlers.ConversationHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *ConversationRoutes {
	return &ConversationRoutes{
		path:      "/conversations",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		tenant:    tenant,
	}
}

func (r *ConversationRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Get("/:id/messages", r.handler.ListMessages)
	protected.Post("/:id/messages", r.handler.SendMessage)
}
