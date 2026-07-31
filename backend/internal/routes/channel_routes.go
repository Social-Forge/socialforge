package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type ChannelRoutes struct {
	path      string
	handler   *handlers.ChannelHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewChannelRoutes(
	handler *handlers.ChannelHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
	tenant *middlewares.TenantMiddleware,
) *ChannelRoutes {
	return &ChannelRoutes{
		path:      "/channels",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		tenant:    tenant,
	}
}

func (r *ChannelRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	protected := route.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/", r.handler.List)
	protected.Post("/", r.handler.Create)
	protected.Get("/:id", r.handler.Get)
	protected.Put("/:id", r.handler.Update)
	protected.Delete("/:id", r.handler.Delete)
}
