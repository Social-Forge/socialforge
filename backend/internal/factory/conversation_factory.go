package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type ConversationFactory struct {
	convSvc  *services.ConversationService
	outbound *services.OutboundService
	handler  *handlers.ConversationHandler
	routes   *routes.ConversationRoutes
}

func NewConversationFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *ConversationFactory {
	convSvc := services.NewConversationService(cont.ConversationRepo, cont.MessageRepo, cont.Logger)
	outbound := services.NewOutboundService(
		cont.ConversationRepo,
		cont.MessageRepo,
		cont.MessageOutboxRepo,
		cont.CentrifugoClient,
		cont.RabbitMQ,
		cont.Logger,
	)
	handler := handlers.NewConversationHandler(mw.ContextMiddleware, convSvc, outbound, cont.Logger)
	return &ConversationFactory{
		convSvc:  convSvc,
		outbound: outbound,
		handler:  handler,
		routes: routes.NewConversationRoutes(
			handler,
			mw.ContextMiddleware,
			mw.AuthMiddleware,
			mw.TenantMiddleware,
		),
	}
}

func (f *ConversationFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
