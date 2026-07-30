package redisclient

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/metrics"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	RedisStorage *RedisClient
	redisOnce    sync.Once
)

type RedisClient struct {
	client  *redis.Client
	metrics *metrics.RedisMetrics
	config  *config.RedisConfig
	isUp    bool
	mu      sync.RWMutex
}

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig, metrics *metrics.RedisMetrics) (*RedisClient, error) {
	var initErr error
	redisOnce.Do(func() {
		if ctx == nil {
			ctx = context.Background()
		}
		subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
		defer cancel()

		client := redis.NewClient(&redis.Options{
			Addr:            cfg.GetAddr(),
			Password:        cfg.Password,
			DB:              cfg.DB,
			PoolSize:        100, // Connection pool size
			MinIdleConns:    10,  // Minimum idle connections
			ConnMaxIdleTime: 10 * time.Minute,
			MaxRetries:      3,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,
		})

		if _, err := client.Ping(subCtx).Result(); err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %w", err)
			config.Logger.Error("Redis connection failed",
				zap.String("addr", cfg.GetAddr()),
				zap.Error(err))
			return
		}

		RedisStorage = &RedisClient{
			client:  client,
			metrics: metrics,
			isUp:    true,
			config:  cfg,
		}

		config.Logger.Info("✅ Redis connected successfully",
			zap.String("addr", cfg.GetAddr()),
			zap.Int("db", cfg.DB),
			zap.Int("pool_size", 100),
		)
	})
	if initErr != nil {
		return nil, initErr
	}
	return RedisStorage, nil
}
func GetRedis() (*RedisClient, error) {
	if RedisStorage == nil {
		return nil, errors.New("redis not initialized: call NewRedisClient first")
	}
	return RedisStorage, nil
}
func (r *RedisClient) Client() *redis.Client {
	return r.client
}
func (r *RedisClient) IsUp() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.isUp
}
func (r *RedisClient) IsClosed() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return !r.isUp
}
func (r *RedisClient) Ping(ctx context.Context) (string, error) {
	ctx, cancel := contextpool.WithTimeoutFallback(ctx, 5*time.Second)
	defer cancel()

	var res string
	err := r.withMetrics(ctx, "ping", func() error {
		var err error
		if res, err = r.client.Ping(ctx).Result(); err != nil {
			r.mu.Lock()
			r.isUp = false
			r.mu.Unlock()
			return err
		}
		r.mu.Lock()
		r.isUp = true
		r.mu.Unlock()
		return nil
	})
	if err != nil {
		return "", err
	}

	return res, nil
}
func (rc *RedisClient) ResetWithContext(ctx context.Context) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "reset", func() error {
		return rc.client.FlushDB(ctx).Err()
	})
}
