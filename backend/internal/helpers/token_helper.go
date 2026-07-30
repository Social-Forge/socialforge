package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	redisclient "github/socialforge/internal/infra/redis-client"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TokenHelper struct {
	client *redisclient.RedisClient
}

func NewTokenHelper(client *redisclient.RedisClient) *TokenHelper {
	return &TokenHelper{client: client}
}
func (ts *TokenHelper) StoreSessionCSRF(ctx context.Context, token string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("csrf:%s", token)
	return ts.client.SetAny(subCtx, keys, token, 1*time.Minute)
}
func (ts *TokenHelper) GetCSRFBySession(ctx context.Context, token string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("csrf:%s", token)
	return ts.client.GetString(subCtx, keys)
}
func (ts *TokenHelper) ClearCSRFBySession(ctx context.Context, token string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("csrf:%s", token)
	return ts.client.DeleteCache(subCtx, keys)
}
func (ts *TokenHelper) SetSessionToken(ctx context.Context, metadata *entity.RedisSessionData, expiry time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("session:%s", metadata.SessionID)
	metadata.LastAccessed = time.Now().Unix()

	config.Logger.Info("💾 [SET] Storing session in Redis",
		zap.String("key", keys),
		zap.String("session_id", metadata.SessionID),
		zap.String("user_id", metadata.UserID.String()),
		zap.Duration("expiry", expiry))

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return ts.client.Setbyte(subCtx, keys, metadataJSON, expiry)
}
func (ts *TokenHelper) GetSessionTokenMetadata(ctx context.Context, sessionID string) (*entity.RedisSessionData, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("session:%s", sessionID)

	config.Logger.Info("🔍 [GET] Looking up session in Redis",
		zap.String("key", keys),
		zap.String("session_id", sessionID))

	data, err := ts.client.GetByte(subCtx, keys)
	if err != nil {
		config.Logger.Error("❌ [GET FAILED] Redis lookup failed",
			zap.String("key", keys),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get session metadata from cache storage: %w", err)
	}

	var metadata entity.RedisSessionData
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		config.Logger.Error("❌ [GET FAILED] Unmarshal failed",
			zap.String("key", keys),
			zap.String("raw_data", string(data)),
			zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	config.Logger.Info("✅ [GET SUCCESS] Session found",
		zap.String("key", keys),
		zap.String("user_id", metadata.UserID.String()))

	return &metadata, nil
}
func (ts *TokenHelper) IsSessionTokenRevoked(ctx context.Context, sessionID string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("revoked_token:%s", sessionID)
	return ts.client.GetBool(subCtx, keys)
}
func (ts *TokenHelper) RevokeSessionToken(ctx context.Context, sessionID string, remainingTTL time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("revoked_token:%s", sessionID)
	return ts.client.SetAny(subCtx, keys, "revoked", remainingTTL)
}
func (ts *TokenHelper) ClearSessionToken(ctx context.Context, sessionID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("session:%s", sessionID)
	return ts.client.DeleteCache(subCtx, keys)
}
func (ts *TokenHelper) DeleteAllExceptCurrent(ctx context.Context, currentSessionID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	currentMeta, err := ts.GetSessionTokenMetadata(subCtx, currentSessionID)
	if err != nil {
		return fmt.Errorf("failed to get current session metadata: %w", err)
	}

	var deletedCount int
	var cursor uint64
	var lastErr error
	pattern := "session:*"

	for {
		keys, nextCursor, err := ts.client.Client().Scan(subCtx, cursor, pattern, 100).Result()
		if err != nil {
			lastErr = fmt.Errorf("scan failed: %w", err)
			break
		}
		for _, key := range keys {
			if key == fmt.Sprintf("session:%s", currentSessionID) {
				continue
			}

			sessionID := strings.TrimPrefix(key, "session:")
			meta, err := ts.GetSessionTokenMetadata(subCtx, sessionID)

			if err != nil {
				// Skip jika data tidak valid atau sudah expired
				if errors.Is(err, redis.Nil) {
					continue
				}
				config.Logger.Warn("Failed to get session metadata",
					zap.String("key", key),
					zap.Error(err))
				continue
			}

			if meta.UserID == currentMeta.UserID {
				if err := ts.client.DeleteCache(subCtx, key); err != nil {
					lastErr = fmt.Errorf("failed to delete session %s: %w", key, err)
					continue
				}
				deletedCount++

				revokeKey := fmt.Sprintf("revoked_token:%s", sessionID)
				ts.client.DeleteCache(subCtx, revokeKey)
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	config.Logger.Info("Deleted old sessions",
		zap.Int("deleted_count", deletedCount),
		zap.String("current_session", currentSessionID),
		zap.String("user_id", currentMeta.UserID.String()))

	return lastErr

}
func (ts *TokenHelper) ClearAllSessionToken(ctx context.Context, sessionID, csrfToken string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	csrfKey := fmt.Sprintf("csrf:%s", csrfToken)
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	revokedKey := fmt.Sprintf("revoked_token:%s", sessionID)

	_, err := ts.client.Client().Pipelined(subCtx, func(pipe redis.Pipeliner) error {
		pipe.Del(subCtx, csrfKey)
		pipe.Del(subCtx, sessionKey)
		pipe.Del(subCtx, revokedKey)
		return nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("redis operation timed out: %w", err)
		}
		return fmt.Errorf("failed to clear session tokens: %w", err)
	}

	return nil
}
func (ts *TokenHelper) UpdateLastAccessed(ctx context.Context, sessionID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	// Get current session data
	metadata, err := ts.GetSessionTokenMetadata(subCtx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session data: %w", err)
	}

	// Update last accessed
	metadata.LastAccessed = time.Now().Unix()

	// Get remaining TTL
	key := fmt.Sprintf("session:%s", sessionID)
	ttl, err := ts.client.Client().TTL(subCtx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get TTL: %w", err)
	}

	if ttl > 0 {
		return ts.SetSessionToken(ctx, metadata, ttl)
	}

	return nil
}
func (ts *TokenHelper) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*entity.RedisSessionData, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	var sessions []*entity.RedisSessionData
	var cursor uint64
	pattern := "session:*"

	for {
		keys, nextCursor, err := ts.client.Client().Scan(subCtx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		for _, key := range keys {
			sessionID := strings.TrimPrefix(key, "session:")
			metadata, err := ts.GetSessionTokenMetadata(subCtx, sessionID)
			if err != nil {
				continue // Skip invalid sessions
			}

			if metadata.UserID == userID {
				sessions = append(sessions, metadata)
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return sessions, nil
}
func (ts *TokenHelper) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	sessions, err := ts.GetUserSessions(subCtx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := ts.ClearSessionToken(subCtx, session.SessionID); err != nil {
			config.Logger.Warn("Failed to delete session",
				zap.String("session_id", session.SessionID),
				zap.Error(err))
		}

		revokeKey := fmt.Sprintf("revoked_token:%s", session.SessionID)
		ts.client.DeleteCache(subCtx, revokeKey)
	}

	config.Logger.Info("Deleted all user sessions",
		zap.String("user_id", userID.String()),
		zap.Int("session_count", len(sessions)))

	return nil
}
