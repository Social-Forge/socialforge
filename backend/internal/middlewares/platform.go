package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type PlatformMiddleware struct {
	ctxinject *ContextMiddleware
	logger    *zap.Logger
}

func NewPlatformMiddleware(ctxinject *ContextMiddleware, logger *zap.Logger) *PlatformMiddleware {
	return &PlatformMiddleware{
		ctxinject: ctxinject,
		logger:    logger,
	}
}
func (m *PlatformMiddleware) Setup() fiber.Handler {
	return func(c fiber.Ctx) error {
		platform := c.Get("X-Platform")
		if len(platform) == 0 {
			platform = "browser"
		}
		c.Locals("platform", platform)
		return c.Next()
	}
}
