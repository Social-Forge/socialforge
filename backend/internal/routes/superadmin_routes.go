package routes

import (
	"github/socialforge/internal/entity"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type SuperadminRoutes struct {
	path      string
	handler   *handlers.SuperadminHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
}

func NewSuperadminRoutes(
	handler *handlers.SuperadminHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
) *SuperadminRoutes {
	return &SuperadminRoutes{
		path:      "/superadmin",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
	}
}

func (r *SuperadminRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	// Superadmin routes are cross-tenant: JWT auth + superadmin role guard,
	// no TenantGuard (not scoped to a single tenant).
	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.auth.RequireRoles(entity.RoleSuperAdmin))

	protected.Get("/tenants", r.handler.ListTenants)
}
