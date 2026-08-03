package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type WebchatRoutes struct {
	path      string
	handler   *handlers.WebchatHandler
	ctxinject *middlewares.ContextMiddleware
}

func NewWebchatRoutes(handler *handlers.WebchatHandler, ctxinject *middlewares.ContextMiddleware) *WebchatRoutes {
	return &WebchatRoutes{path: "/webchat", handler: handler, ctxinject: ctxinject}
}

func (r *WebchatRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)
	// Public: embedded widget endpoints.
	route.Get("/:channelId/session", r.handler.Session)
	route.Post("/:channelId/messages", r.handler.Message)
}
