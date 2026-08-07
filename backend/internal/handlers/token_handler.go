package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenHandler struct {
	ctxinject       *middlewares.ContextMiddleware
	tokenService    *services.TokenService
	centrifugoSecret string
	centrifugoWsURL  string
	logger          *zap.Logger
}

func NewTokenHandler(
	ctxinject *middlewares.ContextMiddleware,
	tokenService *services.TokenService,
	centrifugoSecret string,
	centrifugoWsURL string,
	logger *zap.Logger,
) *TokenHandler {
	return &TokenHandler{
		ctxinject:       ctxinject,
		tokenService:    tokenService,
		centrifugoSecret: centrifugoSecret,
		centrifugoWsURL:  centrifugoWsURL,
		logger:          logger,
	}
}

// GetCentrifugoToken issues a short-lived Centrifugo connection JWT (HS256,
// signed with the shared token_hmac_secret_key) for the authenticated user. The
// client subscribes to tenant/conversation/user channels — those namespaces have
// allow_subscribe_for_client=true, so no per-channel tokens are needed.
func (h *TokenHandler) GetCentrifugoToken(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return helpers.Respond(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}
	if h.centrifugoSecret == "" {
		return helpers.Respond(c, fiber.StatusInternalServerError, "Realtime not configured", nil)
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.centrifugoSecret))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, "Failed to sign realtime token", nil)
	}
	tenantID, _ := c.Locals("tenant_id").(string)
	return helpers.Respond(c, fiber.StatusOK, "Realtime token generated", fiber.Map{
		"token":     token,
		"ws_url":    h.centrifugoWsURL,
		"tenant_id": tenantID,
	})
}
func (h *TokenHandler) GetCSRFToken(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	token, err := h.tokenService.StoreCSRFToken(ctx)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "CSRF token generated", fiber.Map{"csrf_token": token})
}
