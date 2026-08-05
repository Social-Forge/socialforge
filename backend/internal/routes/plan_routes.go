package routes

import (
	"github/socialforge/internal/entity"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type PlanRoutes struct {
	path      string
	handler   *handlers.PlanHandler
	ctxinject *middlewares.ContextMiddleware
	auth      *middlewares.AuthMiddleware
}

func NewPlanRoutes(
	handler *handlers.PlanHandler,
	ctxinject *middlewares.ContextMiddleware,
	auth *middlewares.AuthMiddleware,
) *PlanRoutes {
	return &PlanRoutes{path: "/plans", handler: handler, ctxinject: ctxinject, auth: auth}
}

func (r *PlanRoutes) RegisterRoutes(parent fiber.Router) {
	route := parent.Group(r.path)

	// Public pricing catalog.
	route.Get("/", r.handler.List)
	route.Get("/:code", r.handler.Get)

	// Superadmin plan management (cross-tenant).
	admin := route.Group("/admin")
	admin.Use(r.auth.JWTAuth(), r.auth.RequireRoles(entity.RoleSuperAdmin))
	admin.Get("/", r.handler.ListAll)
	admin.Post("/", r.handler.Create)
	admin.Put("/:id", r.handler.Update)
	admin.Delete("/:id", r.handler.Delete)
}
