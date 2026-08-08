package middlewares

import (
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TenantMiddleware struct {
	notifier     config.Notifier
	ctxinject    *ContextMiddleware
	logger       *zap.Logger
	tenantHelper *helpers.TenantHelper
}

func NewTenantMiddleware(notifier config.Notifier, ctxinject *ContextMiddleware, logger *zap.Logger, tenantHelper *helpers.TenantHelper) *TenantMiddleware {
	return &TenantMiddleware{
		notifier:     notifier,
		ctxinject:    ctxinject,
		logger:       logger,
		tenantHelper: tenantHelper,
	}
}
func (m *TenantMiddleware) TenantGuard() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := m.ctxinject.From(c)

		subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
		defer cancel()

		tenantID := c.Locals("tenant_id").(string)

		if tenantID == "" {
			m.logger.Error("Tenant ID is required")
			return helpers.Respond(c, fiber.StatusBadRequest, "Unauthorized, tenant ID is required", nil)
		}

		tenantUUID, err := uuid.Parse(tenantID)
		if err != nil {
			m.logger.Error("Invalid tenant ID format", zap.Error(err))
			return helpers.Respond(c, fiber.StatusBadRequest, "Invalid tenant ID format", nil)
		}

		tenant, err := m.tenantHelper.GetCacheTenantByID(subCtx, tenantUUID)
		if err != nil || tenant == nil {
			m.logger.Error("Failed to get tenant by ID", zap.Error(err))
			return helpers.Respond(c, fiber.StatusInternalServerError, "Internal server error, tenant not registered", nil)
		}
		if !tenant.IsActive {
			m.logger.Error("Tenant is inactive")
			return helpers.Respond(c, fiber.StatusForbidden, "Tenant is inactive", nil)
		}
		// An active subscription (paid or free plan) is always valid regardless of
		// trial_ends_at. The trial window only grants access when the subscription
		// itself is not active (e.g. a brand-new tenant still trialing).
		if tenant.SubscriptionStatus != entity.StatusActive {
			onTrial := tenant.TrialEndsAt.Valid && tenant.TrialEndsAt.Time.After(time.Now())
			if !onTrial {
				m.logger.Error("Tenant subscription is not active and trial expired",
					zap.String("status", tenant.SubscriptionStatus))
				return helpers.Respond(c, fiber.StatusForbidden, "Tenant subscription is not active", nil)
			}
		}

		// Propagate the tenant id into the request context so that
		// tenant-scoped repositories can bind Postgres RLS
		// (set_config('app.current_tenant', ...)) via RunInTenantTx.
		c.Locals(ctxKey, repository.WithTenantID(ctx, tenantUUID))
		return c.Next()
	}
}
