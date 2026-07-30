package redisclient

import (
	"context"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/infra/contextpool"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (rc *RedisClient) withMetrics(ctx context.Context, cmd string, fn func() error) error {
	_ = ctx

	start := time.Now()
	err := fn()
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil && err != redis.Nil {
		status = "error"
	}

	defer func() {
		if r := recover(); r != nil {
			config.Logger.Error("❌ Redis metric panic", zap.String("cmd", cmd), zap.Any("recover", r))
		}
	}()

	// ✅ Command count metric
	if rc.metrics.CommandCount != nil {
		rc.metrics.CommandCount.With(prometheus.Labels{
			"command":     cmd,
			"environment": "Social Forge",
			"instance":    rc.config.GetInstance(),
			"status":      status,
		}).Inc()
	}

	// ✅ Duration metric
	if rc.metrics.Duration != nil {
		rc.metrics.Duration.With(prometheus.Labels{
			"command":     cmd,
			"environment": "Social Forge",
			"instance":    rc.config.GetInstance(),
			"status":      status,
		}).Observe(duration)
	}

	// ✅ Error count metric (only if actual error)
	if status == "error" && rc.metrics.ErrorCount != nil {
		rc.metrics.ErrorCount.With(prometheus.Labels{
			"command":  cmd,
			"instance": rc.config.GetInstance(),
			"error":    err.Error(), // pastikan string
		}).Inc()
	}

	return err
}

func (rc *RedisClient) Get(key string) ([]byte, error) {
	ctx, cancel := contextpool.WithTimeoutFallback(context.Background(), 5*time.Second)
	defer cancel()

	var val []byte
	err := rc.withMetrics(ctx, "get", func() error {

		var err error
		val, err = rc.client.Get(ctx, key).Bytes()
		if err == redis.Nil {
			return nil
		}
		return err
	})

	if err != nil {
		config.Logger.Warn("Failed to get Redis value",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, err
	}

	return val, nil
}
func (rc *RedisClient) GetByte(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var val []byte
	err := rc.withMetrics(ctx, "get", func() error {

		var err error
		val, err = rc.client.Get(ctx, key).Bytes()
		if err == redis.Nil {
			return nil
		}
		return err
	})

	if err != nil {
		config.Logger.Warn("Failed to get Redis value",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, err
	}

	return val, nil
}
func (rc *RedisClient) GetInt(ctx context.Context, key string) (int, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var value int
	err := rc.withMetrics(ctx, "get_count", func() error {
		cmd := rc.client.Get(ctx, key)
		if cmd.Err() == redis.Nil {
			return nil // Key doesn't exist is not an error case
		}

		var err error
		value, err = cmd.Int()
		return err
	})

	if err != nil {
		if err != redis.Nil {
			config.Logger.Warn("Failed to get bool from Redis", zap.Error(err))
		}
		return value, err
	}
	return value, nil
}
func (rc *RedisClient) GetIntWithDefault(ctx context.Context, key string, defaultValue int) (int, error) {
	val, err := rc.GetInt(ctx, key)
	if err == redis.Nil {
		return defaultValue, nil
	}
	return val, err
}
func (rc *RedisClient) GetString(ctx context.Context, key string) (string, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var val string
	err := rc.withMetrics(ctx, "get_string", func() error {

		var err error
		val, err = rc.client.Get(ctx, key).Result()
		if err == redis.Nil {
			return nil
		}
		return err
	})

	if err != nil {
		if err != redis.Nil {
			config.Logger.Warn("Failed to get bool from Redis", zap.Error(err))
		}
		return val, err
	}

	return val, nil
}
func (rc *RedisClient) GetStringWithDefault(ctx context.Context, key string, defaultValue string) (string, error) {
	val, err := rc.GetString(ctx, key)
	if err == redis.Nil {
		return defaultValue, nil
	}
	return val, err
}
func (rc *RedisClient) GetBool(ctx context.Context, key string) (bool, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var val bool
	err := rc.withMetrics(ctx, "get_bool", func() error {
		strVal, err := rc.client.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}

		val, err = strconv.ParseBool(strVal)
		if err != nil {
			return fmt.Errorf("invalid bool value for key %s: %v", key, err)
		}

		return nil
	})

	if err != nil {
		if err != redis.Nil {
			config.Logger.Warn("Failed to get bool from Redis", zap.Error(err))
		}
		return false, err
	}

	return val, nil
}
func (rc *RedisClient) GetBoolWithDefault(ctx context.Context, key string, defaultValue bool) (bool, error) {
	val, err := rc.GetBool(ctx, key)
	if err == redis.Nil {
		return defaultValue, nil
	}
	return val, err
}

func (rc *RedisClient) Set(key string, val []byte, exp time.Duration) error {
	ctx, cancel := contextpool.WithTimeoutFallback(context.Background(), 5*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "set", func() error {
		err := rc.client.Set(ctx, key, val, exp).Err()
		if err != nil {
			return fmt.Errorf("failed to set Redis key %s: %w", key, err)
		}
		config.Logger.Debug("Successfully set Redis value (bytes)",
			zap.String("key", key),
			zap.Int("value_size", len(val)),
		)
		return err
	})
}
func (rc *RedisClient) SetAny(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "set_any", func() error {
		err := rc.client.Set(ctx, key, val, exp).Err()
		if err != nil {
			return fmt.Errorf("failed to set Redis key %s: %w", key, err)
		}
		return nil
	})
}
func (rc *RedisClient) Setbyte(ctx context.Context, key string, val []byte, exp time.Duration) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "set_byte", func() error {
		err := rc.client.Set(ctx, key, val, exp).Err()
		if err != nil {
			return fmt.Errorf("failed to set Redis key %s: %w", key, err)
		}
		return nil
	})
}
func (rc *RedisClient) SetBool(ctx context.Context, key string, val bool, exp time.Duration) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "set_bool", func() error {
		err := rc.client.Set(ctx, key, val, exp).Err()
		if err != nil {
			return fmt.Errorf("failed to set bool Redis key %s: %w", key, err)
		}

		return nil
	})
}

func (rc *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var incr int64
	err := rc.withMetrics(ctx, "increment", func() error {
		var err error
		incr, err = rc.client.Incr(ctx, key).Result()
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		return 0, err
	}
	return incr, nil
}
func (rc *RedisClient) Scan(ctx context.Context, cursor uint64, match, keys string, count int64) ([]string, uint64, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var nextKeys []string
	var nextCursor uint64
	err := rc.withMetrics(ctx, "token_cleanup", func() error {
		var err error
		nextKeys, nextCursor, err = rc.client.Scan(ctx, cursor, keys, count).Result()
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		return nil, 0, err
	}
	return nextKeys, nextCursor, nil
}
func (rc *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	err := rc.withMetrics(ctx, "key_expire", func() error {
		expired, err := rc.client.Expire(ctx, key, expiration).Result()
		if err != nil {
			return fmt.Errorf("redis expire failed for key %s: %w", key, err)
		}
		if !expired {
			return fmt.Errorf("redis expire not applied for key %s (key may not exist)", key)
		}

		return nil
	})

	if err != nil {
		config.Logger.Warn("Failed to set key expiration", zap.Error(err))
		return err
	}

	return nil
}
func (rc *RedisClient) Exists(ctx context.Context, key string) (int64, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var exists int64
	err := rc.withMetrics(ctx, "exists", func() error {
		var err error
		exists, err = rc.client.Exists(ctx, key).Result()
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		return 0, err
	}
	return exists, err
}
func (rc *RedisClient) Delete(key string) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(context.Background(), 5*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "del", func() error {
		err := rc.client.Del(ctx, key).Err()
		if err != nil {
			return err
		}
		return nil
	})
}
func (rc *RedisClient) DeleteCache(ctx context.Context, key string) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "del", func() error {
		err := rc.client.Del(ctx, key).Err()
		if err != nil {
			return err
		}
		return nil
	})
}
func (rc *RedisClient) Allow(ctx context.Context, key string, limit int, window time.Duration) bool {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var allowed bool
	_ = rc.withMetrics(ctx, "rate_limit", func() error {

		pipe := rc.client.TxPipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)

		_, err := pipe.Exec(ctx)
		if err != nil {
			allowed = false
			return err
		}

		allowed = incr.Val() <= int64(limit)
		return nil
	})
	return allowed
}
func (rc *RedisClient) Pipeline(ctx context.Context) (redis.Pipeliner, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var pipeline redis.Pipeliner
	err := rc.withMetrics(ctx, "pipeline", func() error {
		var errs error
		pipeline = rc.client.Pipeline()
		return errs
	})

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}
func (rc *RedisClient) TxPipeline(ctx context.Context) (redis.Pipeliner, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var pipeline redis.Pipeliner
	err := rc.withMetrics(ctx, "tx_pipeline", func() error {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var errs error
		pipeline = rc.client.TxPipeline()
		return errs
	})

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}
func (rc *RedisClient) HSet(ctx context.Context, key string, values interface{}) (*redis.IntCmd, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var cmd *redis.IntCmd
	err := rc.withMetrics(ctx, "hset", func() error {

		var errSet error
		cmd = rc.client.HSet(ctx, key, values)
		return errSet
	})

	if err != nil {
		return nil, err
	}
	return cmd, nil
}
func (rc *RedisClient) Subscribe(ctx context.Context, key string) (*redis.PubSub, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var cmd *redis.PubSub
	err := rc.withMetrics(ctx, "subscribe", func() error {
		var errSubscribe error
		cmd = rc.client.Subscribe(ctx, key)
		return errSubscribe
	})

	if err != nil {
		return nil, err
	}
	return cmd, nil
}
func (rc *RedisClient) Publish(ctx context.Context, channel string, message interface{}) (*redis.IntCmd, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var cmd *redis.IntCmd
	err := rc.withMetrics(ctx, "publish", func() error {
		var errPublish error
		cmd = rc.client.Publish(ctx, channel, message)
		return errPublish
	})

	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (rc *RedisClient) Reset() error {
	ctx, cancel := contextpool.WithTimeoutFallback(context.Background(), 5*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "reset", func() error {
		// Implementasi reset sesuai kebutuhan
		config.Logger.Warn("Reset operation called but not implemented")
		return nil
	})
}
func (rc *RedisClient) Restore(ctx context.Context, key string, value string, ttl time.Duration) (string, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var result string
	err := rc.withMetrics(ctx, "restore", func() error {
		var err error
		result, err = rc.client.Restore(ctx, key, ttl, value).Result()
		if err != nil {
			return fmt.Errorf("failed to restore data: %w", err)
		}
		return err
	})

	if err != nil {
		return "", fmt.Errorf("failed to restore data: %w", err)
	}
	return result, nil

}
func (rc *RedisClient) Close() error {
	if rc == nil {
		return nil
	}

	if rc.client == nil {
		rc.isUp = false
		return nil
	}

	ctx, cancel := contextpool.WithTimeoutFallback(context.Background(), 5*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "close", func() error {
		rc.mu.Lock()
		defer rc.mu.Unlock()

		if rc.client != nil {
			err := rc.client.Close()
			if err != nil {
				config.Logger.Error("failed to close Redis connection", zap.Error(err))
				rc.isUp = false
				return fmt.Errorf("redis shutdown error: %w", err)
			}
			rc.client = nil
			rc.isUp = false
		}
		return nil
	})
}
func (rc *RedisClient) CloseClient(ctx context.Context) error {
	// Penanganan nil pointer
	if rc == nil {
		return nil
	}

	if rc.client == nil {
		rc.isUp = false
		return nil
	}

	ctx, cancel := contextpool.WithTimeoutFallback(ctx, 5*time.Second)
	defer cancel()

	return rc.withMetrics(ctx, "close", func() error {
		err := rc.client.Close()
		if err != nil {
			config.Logger.Error("failed to close Redis connection", zap.Error(err))
			rc.isUp = false
			return fmt.Errorf("redis shutdown error: %w", err)
		}
		rc.client = nil
		rc.isUp = false
		return nil
	})
}
