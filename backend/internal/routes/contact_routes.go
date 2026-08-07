package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type ContactRoutes struct {
	path      string
	handler   *handlers.ContactHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewContactRoutes(
	handler *handlers.ContactHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *ContactRoutes {
	return &ContactRoutes{path: "/contacts", handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *ContactRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)
	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Post("/:id/block", r.handler.Block)
	protected.Post("/:id/unblock", r.handler.Unblock)
	protected.Delete("/:id", r.handler.Delete)
}
