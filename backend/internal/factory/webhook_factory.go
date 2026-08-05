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
	outbound := services.NewOutboundService(
		cont.ConversationRepo,
		cont.MessageRepo,
		cont.MessageOutboxRepo,
		cont.ChannelRepo,
		cont.ContactRepo,
		cont.CentrifugoClient,
		cont.RabbitMQ,
		services.BuildSenders(cont.Config),
		cont.Logger,
	)
	aiReply := services.NewAIReplyService(
		cont.AIAgentRepo,
		cont.MessageRepo,
		cont.ConversationRepo,
		cont.AIKnowledgeRepo,
		cont.AIPlaybookRepo,
		cont.AIAssetRepo,
		cont.AICreditRepo,
		cont.AIClient,
		outbound,
		cont.Logger,
	)
	service := services.NewIngestionService(
		cont.ChannelRepo,
		cont.ContactRepo,
		cont.ConversationRepo,
		cont.MessageRepo,
		cont.WebhookEventRepo,
		cont.AutoResponseRepo,
		outbound,
		aiReply,
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
