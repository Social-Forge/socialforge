package middlewares

import (
	"github/socialforge/config"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

var (
	ErrNotFoundCsrfToken    = "Unauthorized - CSRF token not found"
	ErrMismatchCsrfToken    = "Unauthorized - CSRF mismatch"
	ErrExpiredCsrfToken     = "Unauthorized - Invalid or expired csrf token"
	ErrInvalidCsrfAuth      = "Unauthorized - Invalid authorization header"
	ErrInvalidOriginRequest = "Forbidden - Invalid request origin"
)

type CSRFMiddleware struct {
	notifier     config.Notifier
	ctxinject    *ContextMiddleware
	tokenHelper  *helpers.TokenHelper
	tenantHelper *helpers.TenantHelper
	logger       *zap.Logger
}

func NewCSRFMiddleware(notifier config.Notifier, ctxinject *ContextMiddleware, tokenHelper *helpers.TokenHelper, tenantHelper *helpers.TenantHelper, logger *zap.Logger) *CSRFMiddleware {
	return &CSRFMiddleware{
		notifier:     notifier,
		ctxinject:    ctxinject,
		tokenHelper:  tokenHelper,
		tenantHelper: tenantHelper,
		logger:       logger,
	}
}
func (m *CSRFMiddleware) CSRFProtect() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := m.ctxinject.From(c)
		defer m.ctxinject.LogDuration(ctx, c.Path())()

		provided := c.Get("X-XSRF-TOKEN")
		if provided == "" {
			m.logger.Error("CSRF token not found", zap.String("path", c.Path()))
			return helpers.Respond(c, fiber.StatusUnauthorized, ErrNotFoundCsrfToken, nil)
		}

		if c.Method() == fiber.MethodGet || c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 60*time.Second)
		defer cancel()

		expected, err := m.tokenHelper.GetCSRFBySession(subCtx, provided)
		if err != nil || expected == "" {
			m.logger.Error("CSRF token not found", zap.String("path", c.Path()))
			return helpers.Respond(c, fiber.StatusUnauthorized, ErrNotFoundCsrfToken, nil)
		}

		return c.Next()
	}
}
