package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type UserFactory struct {
	service *services.UserService
	handler *handlers.UserHandler
	routes  *routes.UserRoutes
}

func NewUserFactory(cont *dependencies.Container, mw *MiddlewareFactory) *UserFactory {
	service := services.NewUserService(
		cont.UserRepo,
		cont.RoleRepo,
		cont.TenantRepo,
		cont.DivisionRepo,
		cont.UserTenantRepo,
		cont.UserHelper,
		cont.TokenHelper,
		cont.Logger,
		cont.MinioClient,
	)
	handler := handlers.NewUserHandler(
		mw.ContextMiddleware,
		service,
		mw.RateLimiter,
		cont.Logger,
	)
	return &UserFactory{
		service: service,
		handler: handler,
		routes: routes.NewUserRoutes(
			handler,
			mw.ContextMiddleware,
			mw.RateLimiter,
			mw.AuthMiddleware,
			mw.CSRFMiddleware,
			mw.TenantMiddleware,
		),
	}
}
func (f *UserFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
