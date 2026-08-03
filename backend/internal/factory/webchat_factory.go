package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type WebchatFactory struct {
	handler *handlers.WebchatHandler
	routes  *routes.WebchatRoutes
}

func NewWebchatFactory(cont *dependencies.Container, mw *MiddlewareFactory) *WebchatFactory {
	outbound := services.NewOutboundService(
		cont.ConversationRepo, cont.MessageRepo, cont.MessageOutboxRepo,
		cont.ChannelRepo, cont.ContactRepo, cont.CentrifugoClient, cont.RabbitMQ,
		services.BuildSenders(cont.Config), cont.Logger,
	)
	ingestion := services.NewIngestionService(
		cont.ChannelRepo, cont.ContactRepo, cont.ConversationRepo, cont.MessageRepo,
		cont.WebhookEventRepo, cont.AutoResponseRepo, outbound,
		cont.CentrifugoClient, cont.RabbitMQ, cont.Logger,
	)
	handler := handlers.NewWebchatHandler(mw.ContextMiddleware, ingestion, cont.Config.Centrifugo.WsURl, cont.Logger)
	return &WebchatFactory{
		handler: handler,
		routes:  routes.NewWebchatRoutes(handler, mw.ContextMiddleware),
	}
}

func (f *WebchatFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
