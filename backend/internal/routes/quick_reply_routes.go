package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type QuickReplyRoutes struct {
	path      string
	handler   *handlers.QuickReplyHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewQuickReplyRoutes(
	handler *handlers.QuickReplyHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *QuickReplyRoutes {
	return &QuickReplyRoutes{path: "/quick-replies", handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *QuickReplyRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Post("/", r.handler.Create)
	protected.Put("/:id", r.handler.Update)
	protected.Delete("/:id", r.handler.Delete)
}
