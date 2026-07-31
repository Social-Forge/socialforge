package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type DivisionRoutes struct {
	path      string
	handler   *handlers.DivisionHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewDivisionRoutes(
	handler *handlers.DivisionHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *DivisionRoutes {
	return &DivisionRoutes{
		path:      "/divisions",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		tenant:    tenant,
	}
}

func (r *DivisionRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Post("/", r.handler.Create)
	protected.Get("/", r.handler.List)
	protected.Get("/:id", r.handler.Get)
	protected.Put("/:id", r.handler.Update)
	protected.Delete("/:id", r.handler.Delete)

	protected.Get("/:id/members", r.handler.ListMembers)
	protected.Post("/:id/members", r.handler.AddMember)
	protected.Delete("/:id/members/:userTenantID", r.handler.RemoveMember)
}
