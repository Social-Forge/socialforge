package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type MemberRoutes struct {
	path      string
	handler   *handlers.MemberHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewMemberRoutes(
	handler *handlers.MemberHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *MemberRoutes {
	return &MemberRoutes{
		path:      "/members",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		tenant:    tenant,
	}
}

func (r *MemberRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Post("/", r.handler.Create)
	protected.Put("/:userTenantID", r.handler.Update)
	protected.Delete("/:userTenantID", r.handler.Delete)
}
