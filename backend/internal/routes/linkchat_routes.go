package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type LinkchatRoutes struct {
	path      string
	handler   *handlers.LinkchatHandler
	ctxinject *middlewares.ContextMiddleware
}

func NewLinkchatRoutes(handler *handlers.LinkchatHandler, ctxinject *middlewares.ContextMiddleware) *LinkchatRoutes {
	return &LinkchatRoutes{path: "/link", handler: handler, ctxinject: ctxinject}
}

func (r *LinkchatRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)
	// Public: shared division link landing.
	route.Get("/:token", r.handler.Resolve)
}
