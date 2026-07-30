package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type TenantRoutes struct {
	path      string
	handler   *handlers.TenantHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
	csrf      *middlewares.CSRFMiddleware
}

func NewTenantRoutes(
	handler *handlers.TenantHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
	csrf *middlewares.CSRFMiddleware,
) *TenantRoutes {
	return &TenantRoutes{
		path:      "/tenants",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		tenant:    tenant,
		csrf:      csrf,
	}
}

func (r *TenantRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Post("/logo", r.handler.UpdateLogo)
	protected.Put("/info/:tenantID", r.handler.UpdateInfo)
}
