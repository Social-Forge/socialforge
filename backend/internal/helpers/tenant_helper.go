package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	redisclient "github/socialforge/internal/infra/redis-client"
	"github/socialforge/internal/infra/repository"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	TenantCacheTTL = 24 * time.Hour
	TenantPrefix   = "tenant:"
)

type TenantHelper struct {
	client               *redisclient.RedisClient
	userRepo             repository.UserRepository
	tenantRepo           repository.TenantRepository
	logger               *zap.Logger
	mu                   sync.RWMutex
	allowedTenantID      map[uuid.UUID]struct{}
	allowedTenantIDSlice []uuid.UUID
	lastUpdated          time.Time
	refreshing           bool
	refreshMu            sync.Mutex
}

func NewTenantHelper(client *redisclient.RedisClient, userRepo repository.UserRepository, tenantRepo repository.TenantRepository, logger *zap.Logger) *TenantHelper {
	return &TenantHelper{
		client:     client,
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		logger:     logger,
	}
}

// Middleware
func (h *TenantHelper) StartTenantRefreshSubscribe(ctx context.Context) {
	if h.client == nil || h.client.Client() == nil {
		config.Logger.Fatal("❌ Redis client not initialized")
		return
	}

	go func() {
		pubsub, err := h.client.Subscribe(ctx, "tenant_updated")
		if err != nil {
			config.Logger.Error("Redis subscribe error", zap.Error(err))
			return
		}
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					config.Logger.Info("Redis Pub/Sub closed", zap.Error(err))
					return
				}
				config.Logger.Warn("Redis Pub/Sub error, retrying...", zap.Error(err))
				time.Sleep(3 * time.Second)
				continue
			}

			var payload map[string]any
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err == nil {
				if payloadType, ok := payload["type"].(string); ok && payloadType == "refresh" {
					config.Logger.Info("🔔 Received tenant_updated -> refreshing allowed tenant IDs")
					go h.safeRefresh(context.Background())
				}
			}
		}
	}()
}
func (h *TenantHelper) InitAllowedTenantIDs(ctx context.Context) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 30*time.Second)
	defer cancel()

	tenantIDs, err := h.tenantRepo.GetAllowedTenantIDs(subCtx)
	if err != nil {
		config.Logger.Error("❌ Failed to load tenant IDs", zap.Error(err))
		return err
	}

	h.updateAllowedTenantIDs(tenantIDs)
	config.Logger.Info("✅ All tenant IDs has stored to cache", zap.Any("tenantIDs", tenantIDs))
	return nil
}
func (h *TenantHelper) RefreshAllowedTenantIDs(ctx context.Context) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tenantIDs, err := h.tenantRepo.GetAllowedTenantIDs(subCtx)
	if err != nil {
		config.Logger.Debug("❌ Failed to refresh tenant IDs", zap.Error(err))
		return err
	}
	h.updateAllowedTenantIDs(tenantIDs)
	config.Logger.Info("🔁 Tenant IDs updated from DB", zap.Any("tenantIDs", tenantIDs))
	return nil
}
func (h *TenantHelper) GetAllowedTenantIDs() []uuid.UUID {
	h.mu.RLock()
	result := make([]uuid.UUID, len(h.allowedTenantIDSlice))
	copy(result, h.allowedTenantIDSlice)
	lastUpdate := h.lastUpdated
	h.mu.RUnlock()

	if time.Since(lastUpdate) > time.Hour {
		go h.safeRefresh(context.Background()) // use background context
	}

	return result
}
func (h *TenantHelper) IsTenantAllowed(tenantID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.allowedTenantID[tenantID]
	return ok
}
func (h *TenantHelper) PublishTenantRefreshSignal(ctx context.Context) error {
	msg := map[string]interface{}{
		"type": "refresh",
		"ts":   time.Now().Unix(),
	}
	data, _ := json.Marshal(msg)
	return h.client.Client().Publish(ctx, "tenant_updated", data).Err()
}
func (h *TenantHelper) safeRefresh(ctx context.Context) {
	h.refreshMu.Lock()
	if h.refreshing {
		h.refreshMu.Unlock()
		return
	}
	h.refreshing = true
	h.refreshMu.Unlock()

	defer func() {
		h.refreshMu.Lock()
		h.refreshing = false
		h.refreshMu.Unlock()
	}()

	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	tenantIDs, err := h.tenantRepo.GetAllowedTenantIDs(subCtx)
	if err != nil {
		config.Logger.Error("🔁 Failed to refresh allowed tenant IDs", zap.Error(err))
		return
	}

	h.updateAllowedTenantIDs(tenantIDs)
	config.Logger.Info("✅ Refreshed tenant IDs from DB", zap.Any("tenantIDs", tenantIDs))
}
func (h *TenantHelper) updateAllowedTenantIDs(tenantIDs []uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.allowedTenantID = make(map[uuid.UUID]struct{})
	for _, o := range tenantIDs {
		h.allowedTenantID[o] = struct{}{}
	}

	h.allowedTenantIDSlice = make([]uuid.UUID, 0, len(h.allowedTenantID))
	for tenantID := range h.allowedTenantID {
		h.allowedTenantIDSlice = append(h.allowedTenantIDSlice, tenantID)
	}

	h.lastUpdated = time.Now()
}

// Helper
func (h *TenantHelper) GetCacheTenantByID(ctx context.Context, tenantID uuid.UUID) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	val, err := h.client.GetByte(subCtx, h.tenantKey(tenantID))
	if err == redis.Nil {
		tenantData, errTenant := h.getTenantByID(subCtx, tenantID)
		if errTenant != nil {
			return nil, errTenant
		}
		err = h.setCacheTenant(subCtx, tenantID, tenantData)
		if err != nil {
			return nil, err
		}
		return tenantData, nil
	}
	tenant := new(entity.Tenant)
	if err = json.Unmarshal([]byte(val), &tenant); err != nil {
		if err = h.DeleteCacheTenantByID(subCtx, tenantID); err != nil {
			h.logger.Warn("failed to delete cache tenant by tenant ID", zap.Error(err), zap.Any("tenantID", tenantID))
		}
		tenantData, err := h.getTenantByID(subCtx, tenantID)
		if err != nil {
			return nil, err
		}
		if err = h.setCacheTenant(subCtx, tenantID, tenantData); err != nil {
			return nil, err
		}

		return tenantData, nil
	}

	return tenant, nil
}
func (h *TenantHelper) DeleteCacheTenantByID(ctx context.Context, tenantID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return h.client.DeleteCache(subCtx, h.tenantKey(tenantID))
}
func (h *TenantHelper) tenantKey(tenantID uuid.UUID) string {
	return fmt.Sprintf("%s%s", TenantPrefix, tenantID)
}
func (h *TenantHelper) getTenantByID(ctx context.Context, tenantID uuid.UUID) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tenant, err := h.tenantRepo.FindByID(subCtx, tenantID)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}
func (h *TenantHelper) setCacheTenant(ctx context.Context, userID uuid.UUID, tenant *entity.Tenant) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(tenant)
	if err != nil {
		return fmt.Errorf("failed to marshal tenant: %w", err)
	}
	err = h.client.Setbyte(subCtx, h.tenantKey(userID), jsonData, 10*time.Minute)
	if err != nil {
		return err
	}
	return nil
}
func StructToMap(data interface{}) map[string]interface{} {
	var result map[string]interface{}
	jsonData, _ := json.Marshal(data)
	_ = json.Unmarshal(jsonData, &result)
	return result
}
