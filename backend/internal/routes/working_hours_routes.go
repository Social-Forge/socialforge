package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type WorkingHoursRoutes struct {
	path      string
	handler   *handlers.WorkingHoursHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewWorkingHoursRoutes(
	handler *handlers.WorkingHoursHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *WorkingHoursRoutes {
	return &WorkingHoursRoutes{path: "/working-hours", handler: handler, ctxinject: ctxinject, auth: auth, tenant: tenant}
}

func (r *WorkingHoursRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/:userId", r.handler.List)
	protected.Put("/:userId", r.handler.Set)
}
