package factory

import (
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/middlewares"
)

type MiddlewareFactory struct {
	ContextMiddleware      *middlewares.ContextMiddleware
	Recovery               *middlewares.RecoveryMiddleware
	RateLimiter            *middlewares.RateLimiterMiddleware
	ApiMiddleware          *middlewares.ApiMiddleware
	AuthMiddleware         *middlewares.AuthMiddleware
	TenantMiddleware       *middlewares.TenantMiddleware
	CSRFMiddleware         *middlewares.CSRFMiddleware
	RequireFlagsMiddleware *middlewares.RequireFlagsMiddleware
	PlatformMiddleware     *middlewares.PlatformMiddleware
}

func NewMiddlewareFactory(cont *dependencies.Container) *MiddlewareFactory {
	ctxinject := middlewares.NewContextMiddleware(cont.Logger)
	return &MiddlewareFactory{
		ContextMiddleware:      ctxinject,
		Recovery:               middlewares.NewRecoveryMiddleware(ctxinject, cont.Logger, cont.Notifier),
		RateLimiter:            middlewares.NewRateLimiterMiddleware(ctxinject, cont.RedisClient),
		ApiMiddleware:          middlewares.NewApiMiddleware(),
		AuthMiddleware:         middlewares.NewAuthMiddleware(cont.Notifier, ctxinject, cont.RedisClient, cont.TokenHelper, cont.Logger, cont.Config.JWT.Secret),
		TenantMiddleware:       middlewares.NewTenantMiddleware(cont.Notifier, ctxinject, cont.Logger, cont.TenantHelper),
		CSRFMiddleware:         middlewares.NewCSRFMiddleware(cont.Notifier, ctxinject, cont.TokenHelper, cont.TenantHelper, cont.Logger),
		RequireFlagsMiddleware: middlewares.NewRequireFlagsMiddleware(ctxinject, cont.UserHelper, cont.Logger),
		PlatformMiddleware:     middlewares.NewPlatformMiddleware(ctxinject, cont.Logger),
	}
}
