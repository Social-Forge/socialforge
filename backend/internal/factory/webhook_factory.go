package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/routes"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
)

type WebhookFactory struct {
	service *services.IngestionService
	handler *handlers.WebhookHandler
	routes  *routes.WebhookRoutes
}

func NewWebhookFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *WebhookFactory {
	service := services.NewIngestionService(
		cont.ChannelRepo,
		cont.ContactRepo,
		cont.ConversationRepo,
		cont.MessageRepo,
		cont.WebhookEventRepo,
		cont.CentrifugoClient,
		cont.RabbitMQ,
		cont.Logger,
	)
	handler := handlers.NewWebhookHandler(mw.ContextMiddleware, service, cont.Logger)
	return &WebhookFactory{
		service: service,
		handler: handler,
		routes:  routes.NewWebhookRoutes(handler, mw.ContextMiddleware),
	}
}

func (f *WebhookFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
