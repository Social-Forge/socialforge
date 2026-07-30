package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApiMiddleware struct {
}

func NewApiMiddleware() *ApiMiddleware {
	return &ApiMiddleware{}
}
func (m *ApiMiddleware) SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: nil,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin, Referer, Host, Content-Type, Accept, X-Forwarded-Origin, X-Forwarded-Host, Authorization, X-Client-Platform, X-Package-ID, X-XSRF-TOKEN, X-Xsrf-Token, X-Requested-With, X-Original-Url, X-Forwarded-Referer, X-Real-Host, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, User-Agent, X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, X-2FA-Session, X-Require-Confirm, X-Platform"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "QUERY", "OPTIONS"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Content-Length, X-Request-ID, X-Require-Confirm, X-2FA-Session"},
		MaxAge:           86400,
	})
}
func (m *ApiMiddleware) SetupCompression() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
func (m *ApiMiddleware) SetupLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${ip} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	})
}
func (m *ApiMiddleware) SetupRequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
	})
}
func (m *ApiMiddleware) SetupMetrics(log *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		log.Debug("Request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("request_id", c.GetRespHeader("X-Request-ID")),
		)

		return err
	}
}
