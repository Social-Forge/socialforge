package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type WebhookHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.IngestionService
	logger    *zap.Logger
}

func NewWebhookHandler(ctxinject *middlewares.ContextMiddleware, service *services.IngestionService, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{ctxinject: ctxinject, service: service, logger: logger}
}

// Verify handles the Meta webhook verification handshake:
//
//	GET /api/webhooks/:provider/:channelId?hub.mode=subscribe&hub.verify_token=..&hub.challenge=..
//
// It must echo hub.challenge as plain text when the token matches.
func (h *WebhookHandler) Verify(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	if c.Query("hub.mode") == "subscribe" && h.service.VerifyChallenge(ctx, c.Params("channelId"), c.Query("hub.verify_token")) {
		return c.SendString(c.Query("hub.challenge"))
	}
	return c.SendStatus(fiber.StatusForbidden)
}

// Receive is the public inbound webhook endpoint for all providers:
//
//	POST /api/webhooks/:provider/:channelId
func (h *WebhookHandler) Receive(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	provider := c.Params("provider")
	channelID := c.Params("channelId")
	body := c.Body()

	headers := map[string]string{
		"X-Webhook-Secret":                c.Get("X-Webhook-Secret"),
		"X-Telegram-Bot-Api-Secret-Token": c.Get("X-Telegram-Bot-Api-Secret-Token"),
		"X-Webhook-Hmac":                  c.Get("X-Webhook-Hmac"),
		"X-Hub-Signature-256":             c.Get("X-Hub-Signature-256"),
	}

	if err := h.service.VerifyAndEnqueue(ctx, provider, channelID, headers, body); err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "signature") || strings.Contains(msg, "secret"):
			return helpers.Respond(c, fiber.StatusUnauthorized, "unauthorized webhook", nil)
		case strings.Contains(msg, "not found") || strings.Contains(msg, "does not match") ||
			strings.Contains(msg, "invalid channel") || strings.Contains(msg, "unsupported provider"):
			h.logger.Warn("webhook rejected", zap.String("provider", provider), zap.Error(err))
			return helpers.Respond(c, fiber.StatusBadRequest, msg, nil)
		default:
			// Transient/persistence error — return 500 so the provider retries
			// (dedup on webhook_events + provider_message_id makes retries safe).
			h.logger.Error("webhook processing failed", zap.String("provider", provider), zap.Error(err))
			return helpers.Respond(c, fiber.StatusInternalServerError, "processing failed", nil)
		}
	}
	return helpers.Respond(c, fiber.StatusOK, "ok", nil)
}
