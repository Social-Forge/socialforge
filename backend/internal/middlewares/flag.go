package middlewares

import (
	"github/socialforge/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type RequireFlagsMiddleware struct {
	ctxinject  *ContextMiddleware
	userHelper *helpers.UserHelper
	logger     *zap.Logger
}

func NewRequireFlagsMiddleware(ctxinject *ContextMiddleware, userHelper *helpers.UserHelper, logger *zap.Logger) *RequireFlagsMiddleware {
	return &RequireFlagsMiddleware{
		ctxinject:  ctxinject,
		userHelper: userHelper,
		logger:     logger,
	}
}
func (m *RequireFlagsMiddleware) RequireConfirmedPassword() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := m.ctxinject.From(c)
		defer m.ctxinject.LogDuration(ctx, c.Path())()

		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return helpers.Respond(c, fiber.StatusUnauthorized, helpers.ErrUnauthorizedSession, nil)
		}

		confirmed, err := m.userHelper.IsPasswordConfirmed(ctx, userID)
		if err != nil {
			return helpers.Respond(c, fiber.StatusInternalServerError, "Failed to verify session", nil)
		}
		if !confirmed {
			c.Set("X-Require-Confirm", "true")
			return helpers.Respond(c,
				fiber.StatusAccepted,
				"Password confirmation required",
				fiber.Map{
					"requred_confirm": true,
				})
		}

		return c.Next()
	}
}
func (m *RequireFlagsMiddleware) RequireTwoFactor() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := m.ctxinject.From(c)
		defer m.ctxinject.LogDuration(ctx, c.Path())()

		sessionID := c.Cookies("2fa_session_id") // Ambil 2FA session ID dari cookie
		if sessionID == "" {
			return helpers.Respond(c, fiber.StatusForbidden, "2FA session missing", nil)
		}

		status, err := m.userHelper.Get2FaStatus(ctx, sessionID, "verified_2fa")
		if err != nil || status != "true" {
			return helpers.Respond(c, fiber.StatusInternalServerError, "2FA verification required", nil)
		}

		if status != "true" {
			return helpers.Respond(c, fiber.StatusForbidden, "2FA verification required", nil)
		}

		return c.Next()
	}

}
