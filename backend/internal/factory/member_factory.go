package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type MemberFactory struct {
	service *services.MemberService
	handler *handlers.MemberHandler
	routes  *routes.MemberRoutes
}

func NewMemberFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *MemberFactory {
	service := services.NewMemberService(
		cont.UserRepo,
		cont.RoleRepo,
		cont.TenantRepo,
		cont.Logger,
	)
	handler := handlers.NewMemberHandler(
		mw.ContextMiddleware,
		service,
		cont.Logger,
	)
	return &MemberFactory{
		service: service,
		handler: handler,
		routes: routes.NewMemberRoutes(
			handler,
			mw.ContextMiddleware,
			mw.AuthMiddleware,
			mw.TenantMiddleware,
		),
	}
}

func (f *MemberFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
