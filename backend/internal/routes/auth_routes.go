package routes

import (
	"github/socialforge/internal/handlers"
	"github/socialforge/internal/middlewares"

	"time"

	"github.com/gofiber/fiber/v3"
)

type AuthRoutes struct {
	path      string
	handler   *handlers.AuthHandler
	ctxinject *middlewares.ContextMiddleware
	rateLimit *middlewares.RateLimiterMiddleware
}

func NewAuthRoutes(handler *handlers.AuthHandler, ctxinject *middlewares.ContextMiddleware, rateLimit *middlewares.RateLimiterMiddleware) *AuthRoutes {
	return &AuthRoutes{
		path:      "/auth",
		handler:   handler,
		ctxinject: ctxinject,
		rateLimit: rateLimit,
	}
}
func (r *AuthRoutes) RegisterRoutes(app fiber.Router) {
	router := app.Group(r.path)
	router.Post("/register", r.handler.Register)
	router.Post("/login",
		r.ctxinject.SetTimeout(10*time.Second),
		r.rateLimit.BaseLimiter("login", 5, 5*time.Minute),
		r.rateLimit.ProgressDelay("login"),
		r.rateLimit.BlockLimiter("login", 3, 30*time.Minute),
		r.handler.Login)
	router.Post("/forgot",
		r.ctxinject.SetTimeout(10*time.Second),
		r.rateLimit.BaseLimiter("forgot", 5, 5*time.Minute),
		r.rateLimit.ProgressDelay("forgot"),
		r.rateLimit.BlockLimiter("forgot", 3, 30*time.Minute),
		r.handler.ForgotPassword,
	)
	router.Post("/verify-email", r.handler.VerifyEmail)
	router.Post("/reset-password", r.handler.ResetPassword)
	router.Post("/verify-two-factor",
		r.ctxinject.SetTimeout(10*time.Second),
		r.rateLimit.BaseLimiter("verify2fa", 5, 5*time.Minute),
		r.rateLimit.ProgressDelay("verify2fa"),
		r.rateLimit.BlockLimiter("verify2fa", 3, 30*time.Minute),
		r.handler.VerifyTwoFactor,
	)
	router.Post("/refresh-token", r.handler.RefreshToken)
}
