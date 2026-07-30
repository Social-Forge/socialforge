package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

type UserRoutes struct {
	path      string
	handler   *handlers.UserHandler
	ctxinject *middlewares.ContextMiddleware
	limiter   *middlewares.RateLimiterMiddleware
	auth      *middlewares.AuthMiddleware
	csrf      *middlewares.CSRFMiddleware
	tenant    *middlewares.TenantMiddleware
}

func NewUserRoutes(
	handler *handlers.UserHandler,
	ctxinject *middlewares.ContextMiddleware,
	limiter *middlewares.RateLimiterMiddleware,
	auth *middlewares.AuthMiddleware,
	csrf *middlewares.CSRFMiddleware,
	tenant *middlewares.TenantMiddleware,
) *UserRoutes {
	return &UserRoutes{
		path:      "/user",
		handler:   handler,
		ctxinject: ctxinject,
		limiter:   limiter,
		auth:      auth,
		csrf:      csrf,
		tenant:    tenant,
	}
}
func (r *UserRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	protected := router.Group("/protected")
	protected.Use(r.auth.JWTAuth(), r.tenant.TenantGuard())

	protected.Get("/me", r.handler.GetCurrentUser)
	protected.Post("/logout", r.handler.Logout)

	protected.Post("/avatar", r.handler.ChangeAvatar)
	protected.Put("/profile", r.handler.UpdateProfile)
	protected.Put("/password", r.handler.ChangePassword)
	protected.Post("/two-factor/enable", r.handler.EnableTwoFactor)
	protected.Post("/two-factor/verify", r.handler.VerifyTwoFactor)
}
