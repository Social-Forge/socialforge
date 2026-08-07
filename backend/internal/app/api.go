package app

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/factory"

	"github.com/gofiber/fiber/v3"
)

func RegisterApiRoutes(router fiber.Router, cont *dependencies.Container, mw *factory.MiddlewareFactory) {
	authFactory := factory.NewAuthFactory(cont, mw)
	authFactory.GetRoutes(router)

	userFactory := factory.NewUserFactory(cont, mw)
	userFactory.GetRoutes(router)

	tokenFactory := factory.NewTokenFactory(cont, mw)
	tokenFactory.GetRoutes(router)

	tenantFactory := factory.NewTenantFactory(cont, mw)
	tenantFactory.GetRoutes(router)

	divisionFactory := factory.NewDivisionFactory(cont, mw)
	divisionFactory.GetRoutes(router)

	memberFactory := factory.NewMemberFactory(cont, mw)
	memberFactory.GetRoutes(router)

	superadminFactory := factory.NewSuperadminFactory(cont, mw)
	superadminFactory.GetRoutes(router)

	channelFactory := factory.NewChannelFactory(cont, mw)
	channelFactory.GetRoutes(router)

	webhookFactory := factory.NewWebhookFactory(cont, mw)
	webhookFactory.GetRoutes(router)

	conversationFactory := factory.NewConversationFactory(cont, mw)
	conversationFactory.GetRoutes(router)

	labelFactory := factory.NewLabelFactory(cont, mw)
	labelFactory.GetRoutes(router)

	quickReplyFactory := factory.NewQuickReplyFactory(cont, mw)
	quickReplyFactory.GetRoutes(router)

	workingHoursFactory := factory.NewWorkingHoursFactory(cont, mw)
	workingHoursFactory.GetRoutes(router)

	linkchatFactory := factory.NewLinkchatFactory(cont, mw)
	linkchatFactory.GetRoutes(router)

	webchatFactory := factory.NewWebchatFactory(cont, mw)
	webchatFactory.GetRoutes(router)

	aiAgentFactory := factory.NewAIAgentFactory(cont, mw)
	aiAgentFactory.GetRoutes(router)

	aiResourceFactory := factory.NewAIResourceFactory(cont, mw)
	aiResourceFactory.GetRoutes(router)

	planFactory := factory.NewPlanFactory(cont, mw)
	planFactory.GetRoutes(router)

	billingFactory := factory.NewBillingFactory(cont, mw)
	billingFactory.GetRoutes(router)

	contactFactory := factory.NewContactFactory(cont, mw)
	contactFactory.GetRoutes(router)
}
