package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	ctxKey     = "ctx"
	timeoutKey = "route_timeout"
	minTimeout = 5 * time.Second
	maxTimeout = 120 * time.Second
)

type ContextMiddleware struct {
	logger *zap.Logger
}

func NewContextMiddleware(
	logger *zap.Logger,
) *ContextMiddleware {
	return &ContextMiddleware{
		logger: logger,
	}
}
func (cm *ContextMiddleware) TimeoutContext(defaultTimeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		parentCtx := c.Context()
		if parentCtx == nil {
			parentCtx = context.Background()
		}

		routeTimeout := defaultTimeout
		if t, ok := c.Locals(timeoutKey).(time.Duration); ok {
			routeTimeout = t
		}

		if routeTimeout <= 0 {
			routeTimeout = defaultTimeout
		}
		if routeTimeout < minTimeout {
			routeTimeout = minTimeout
		}
		if routeTimeout > maxTimeout {
			routeTimeout = maxTimeout
		}

		ctx, cancel := context.WithTimeout(parentCtx, routeTimeout)
		defer cancel()

		c.Locals("ctx", ctx)

		if deadline, ok := ctx.Deadline(); ok {
			cm.logger.Debug("Context initialized", zap.String("path", c.Path()), zap.Duration("timeout", routeTimeout), zap.Time("deadline", deadline))
		}

		return c.Next()
	}
}
func (cm *ContextMiddleware) SetTimeout(duration time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals(timeoutKey, duration)
		cm.logger.Debug("The context time has been set plus for the route", zap.String("path", c.Path()), zap.Duration("duration", duration))
		return c.Next()
	}
}
func (cm *ContextMiddleware) From(c fiber.Ctx) context.Context {
	if ctx, ok := c.Locals(ctxKey).(context.Context); ok && ctx != nil {
		return ctx
	}
	return context.Background()
}
func (cm *ContextMiddleware) HandlerContext(c fiber.Ctx) context.Context {
	ctx, ok := c.Locals(ctxKey).(context.Context)
	if !ok || ctx == nil {
		return context.Background()
	}

	return ctx
}
func (cm *ContextMiddleware) LogDuration(ctx context.Context, path string) func() {
	// now := time.Now()
	//if deadline, ok := ctx.Deadline(); ok {
	// return func() {}
	// config.Logger.Info("⏳ Context deadline", zap.Time("deadline", deadline), zap.Duration("remaining", time.Until(deadline)))
	//}

	return func() {
		// cm.logger.Info("⏱ Handler/Query Duration",
		// 	zap.String("path", path),
		// 	zap.Duration("duration", time.Since(now)),
		// )
	}
}
