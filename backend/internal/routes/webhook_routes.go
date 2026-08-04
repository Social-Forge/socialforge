package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type WebhookRoutes struct {
	path      string
	handler   *handlers.WebhookHandler
	ctxinject *middlewares.ContextMiddleware
}

func NewWebhookRoutes(handler *handlers.WebhookHandler, ctxinject *middlewares.ContextMiddleware) *WebhookRoutes {
	return &WebhookRoutes{path: "/webhooks", handler: handler, ctxinject: ctxinject}
}

func (r *WebhookRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	// Public: providers post inbound events here. Authenticity is verified per
	// channel via the webhook secret / provider signature inside the handler.
	route.Get("/:provider/:channelId", r.handler.Verify) // Meta hub.challenge handshake
	route.Post("/:provider/:channelId", r.handler.Receive)
}
