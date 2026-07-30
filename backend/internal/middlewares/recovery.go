package middlewares

import (
	"context"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/helpers"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RecoveryMiddleware struct {
	ctxinject *ContextMiddleware
	logger    *zap.Logger
	notifier  config.Notifier
}

func NewRecoveryMiddleware(
	ctxinject *ContextMiddleware,
	logger *zap.Logger,
	notifier config.Notifier,
) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		ctxinject: ctxinject,
		logger:    logger,
		notifier:  notifier,
	}
}
func (rm *RecoveryMiddleware) NewRecoveryMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				originalCtx := rm.ctxinject.From(c)

				isTimeout := originalCtx.Err() == context.DeadlineExceeded
				isCanceled := originalCtx.Err() == context.Canceled

				rm.logger.Error("⚠️ Panic recovered",
					zap.String("path", c.Path()),
					zap.Any("error", r),
					zap.Bool("is_timeout", isTimeout),
					zap.Bool("is_canceled", isCanceled),
					zap.String("stack", string(debug.Stack())),
				)

				if !isTimeout {
					if rm.notifier != nil {
						rm.notifier.SendAlert(config.AlertRequest{
							Subject: "🚨 Panic Recovered in Recovery Middleware",
							Message: fmt.Sprintf("Path: %s\nError: %v\n\nStack: %s", c.Path(), r, string(debug.Stack())),
							Metadata: map[string]interface{}{ // Stack sudah di Message
								"stack": string(debug.Stack()),
							},
						})
					} else {
						rm.logger.Warn("Notifier or recipient list not configured for panic alerts.")
					}
				}

				status := fiber.StatusInternalServerError
				message := "Internal Server Error"

				if isTimeout {
					status = fiber.StatusGatewayTimeout
					message = "Request Timeout"
				}

				c.Locals(timeoutKey, 10*time.Second)

				_ = helpers.Respond(c, status, message, fiber.Map{
					"request_id":  c.Locals("request_id"),
					"incident_id": uuid.New().String(),
					"is_timeout":  isTimeout,
				})
			}
		}()
		return c.Next()
	}
}
