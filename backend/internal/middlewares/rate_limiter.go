package middlewares

import (
	"context"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	redisclient "github/socialforge/internal/infra/redis-client"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RateLimiterMiddleware struct {
	ctxinjext   *ContextMiddleware
	redisClient *redisclient.RedisClient
}

func NewRateLimiterMiddleware(
	ctxinject *ContextMiddleware,
	redisClient *redisclient.RedisClient,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		ctxinjext:   ctxinject,
		redisClient: redisClient,
	}
}
func (rm *RateLimiterMiddleware) ProgressDelay(key string) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := rm.ctxinjext.From(c)

		attemptsKey := fmt.Sprintf("delay:%s:%s", key, c.Locals("real_ip").(string))

		attempts, err := rm.redisClient.GetInt(ctx, attemptsKey)
		if err != nil && err != redis.Nil {
			config.Logger.Error("Redis error", zap.Error(err))
			return c.Next()
		}

		if attempts >= 3 {
			delay := time.Duration(attempts-2) * time.Second
			time.Sleep(delay)
		}
		return c.Next()
	}
}
func (rm *RateLimiterMiddleware) ResetLimitCounters(c fiber.Ctx) {
	ctx := rm.ctxinjext.From(c)

	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	ip := c.Locals("real_ip").(string)
	patterns := []string{
		fmt.Sprintf("rate:%s:%s", "login", ip),
		fmt.Sprintf("delay:%s:%s", "forgot", ip),
		fmt.Sprintf("block:%s:%s", "confirm_password", ip),
	}

	for _, pattern := range patterns {
		rm.redisClient.DeleteCache(ctx, pattern)
	}
}
func (rm *RateLimiterMiddleware) GlobalRequestLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c fiber.Ctx) string {
			ip := c.Locals("real_ip").(string)
			return fmt.Sprintf("global:%s:%s", ip, c.Method())
		},
		Max:        100,
		Expiration: 1 * time.Minute,
		Storage:    &redisStorageWrapper{client: rm.redisClient},
		LimitReached: func(c fiber.Ctx) error {
			retryAfter := c.GetRespHeader("Retry-After")
			return helpers.Respond(c, fiber.StatusTooManyRequests, "global_rate_limit_exceeded", fiber.Map{
				"retry_after": retryAfter,
			})
		},
	})
}

func (rm *RateLimiterMiddleware) BaseLimiter(key string, max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c fiber.Ctx) string {
			ip := c.Locals("real_ip").(string)
			return fmt.Sprintf("rate:%s:%s", key, ip)
		},
		Storage:      &redisStorageWrapper{client: rm.redisClient},
		Max:          max,
		Expiration:   expiration,
		LimitReached: defaultLimitReachedHandler,
	})
}

func (rm *RateLimiterMiddleware) BlockLimiter(key string, maxAttempts int, blockDuration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c fiber.Ctx) string {
			ip := c.Locals("real_ip").(string)
			return fmt.Sprintf("block:%s:%s", key, ip)
		},
		Storage:      &redisStorageWrapper{client: rm.redisClient},
		Max:          maxAttempts,
		Expiration:   blockDuration,
		LimitReached: blockLimitReachedHandler(blockDuration),
	})
}
func defaultLimitReachedHandler(c fiber.Ctx) error {
	return helpers.Respond(c, fiber.StatusTooManyRequests, "You have reached the request limit. Please try again later.", nil)
}
func blockLimitReachedHandler(blockDuration time.Duration) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		msg := fmt.Sprintf("Too many attempts. Please try again after %v.", blockDuration)
		return helpers.Respond(c, fiber.StatusTooManyRequests, msg, fiber.Map{
			"block_duration": blockDuration,
		})
	}
}

type redisStorageWrapper struct {
	client *redisclient.RedisClient
}

func (w *redisStorageWrapper) Get(key string) ([]byte, error) {
	return w.client.Get(key)
}

func (w *redisStorageWrapper) Set(key string, val []byte, exp time.Duration) error {
	return w.client.Set(key, val, exp)
}

func (w *redisStorageWrapper) Delete(key string) error {
	return w.client.Delete(key)
}

func (w *redisStorageWrapper) Reset() error {
	return w.client.Reset()
}

func (w *redisStorageWrapper) Close() error {
	return w.client.Close()
}

func (w *redisStorageWrapper) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	return w.client.GetByte(ctx, key)
}

func (w *redisStorageWrapper) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error {
	return w.client.Setbyte(ctx, key, val, exp)
}

func (w *redisStorageWrapper) DeleteWithContext(ctx context.Context, key string) error {
	return w.client.DeleteCache(ctx, key)
}

func (w *redisStorageWrapper) ResetWithContext(ctx context.Context) error {
	// Implementasi reset dengan context
	return w.client.Reset() // atau buat method ResetWithContext di RedisClient
}

func (w *redisStorageWrapper) CloseWithContext(ctx context.Context) error {
	return w.client.CloseClient(ctx)
}
