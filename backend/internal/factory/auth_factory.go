package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type AuthFactory struct {
	service *services.AuthService
	handler *handlers.AuthHandler
	routes  *routes.AuthRoutes
}

func NewAuthFactory(container *dependencies.Container, mw *MiddlewareFactory) *AuthFactory {
	service := services.NewAuthService(
		container.UserRepo,
		container.RoleRepo,
		container.PermissionRepo,
		container.SessionRepo,
		container.TenantRepo,
		container.UserTenantRepo,
		container.TokenRepo,
		mw.RateLimiter,
		container.TokenHelper,
		container.AuthHelper,
		container.UserHelper,
		container.Logger,
		container.Config.JWT.Secret,
		container.Config.JWT.ExpireHours,
		container.Config.JWT.RefreshExpireHours)

	handler := handlers.NewAuthHandler(
		mw.ContextMiddleware,
		service,
		mw.RateLimiter,
		container.Logger,
	)
	routes := routes.NewAuthRoutes(handler, mw.ContextMiddleware, mw.RateLimiter)
	return &AuthFactory{
		service: service,
		handler: handler,
		routes:  routes,
	}
}

func (f *AuthFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
