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
	protected.Get("/unread", r.handler.Unread)
	protected.Get("/:id/messages", r.handler.ListMessages)
	protected.Post("/:id/messages", r.handler.SendMessage)

	// Message actions
	protected.Post("/:id/messages/:messageId/pin", r.handler.PinMessage)
	protected.Post("/:id/messages/:messageId/unpin", r.handler.UnpinMessage)
	protected.Put("/:id/messages/:messageId", r.handler.EditMessage)
	protected.Delete("/:id/messages/:messageId", r.handler.DeleteMessage)
	protected.Post("/:id/messages/:messageId/forward", r.handler.ForwardMessage)

	// Conversation actions
	protected.Post("/:id/assign", r.handler.Assign)
	protected.Post("/:id/unassign", r.handler.Unassign)
	protected.Post("/:id/complete", r.handler.Complete)
	protected.Post("/:id/reopen", r.handler.Reopen)
	protected.Post("/:id/read", r.handler.MarkRead)
	protected.Post("/:id/pin", r.handler.Pin)
	protected.Post("/:id/unpin", r.handler.Unpin)
	protected.Post("/:id/archive", r.handler.Archive)
	protected.Post("/:id/unarchive", r.handler.Unarchive)
}
