package helpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github/socialforge/internal/infra/contextpool"
	redisclient "github/socialforge/internal/infra/redis-client"
	"github/socialforge/internal/infra/repository"

	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	UserCacheTTL           = 24 * time.Hour
	UserPrefix             = "user:"
	ConfirmPassPerfix      = "confirm_password:"
	TwoFaCodePrefix        = "2fa:"
	TwoFaStatusPrefix      = "status_2fa_user:"
	TempSecrets            = make(map[string]string)
	BlockedCredential      = "blocked_credential:"
	DelayBlockedCredential = "delay_blocked_credential:"
)

type UserHelper struct {
	client   *redisclient.RedisClient
	userRepo repository.UserRepository
}

func NewUserHelper(client *redisclient.RedisClient, userRepo repository.UserRepository) *UserHelper {
	return &UserHelper{
		client:   client,
		userRepo: userRepo,
	}
}
func (us *UserHelper) SetPasswordConfirmed(ctx context.Context, userID string, status bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s%s", ConfirmPassPerfix, userID)
	err := us.client.SetAny(subCtx, key, status, UserCacheTTL)
	if err != nil {
		return fmt.Errorf("failed to set password confirmation: %w", err)
	}
	return nil
}
func (us *UserHelper) IsPasswordConfirmed(ctx context.Context, userID string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s%s", ConfirmPassPerfix, userID)
	val, err := us.client.GetBool(subCtx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get password confirmation: %w", err)
	}
	return val, nil
}
func (us *UserHelper) ClearPasswordConfirmed(ctx context.Context, userID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s", ConfirmPassPerfix, userID)
	return us.client.DeleteCache(subCtx, keys)
}
func (us *UserHelper) Set2FaStatus(ctx context.Context, sessionID, key, value string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s:%s", TwoFaStatusPrefix, sessionID, key)
	return us.client.SetAny(subCtx, keys, value, 10*time.Minute)
}
func (us *UserHelper) Get2FaStatus(ctx context.Context, sessionID, key string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s:%s", TwoFaStatusPrefix, sessionID, key)
	val, err := us.client.GetString(subCtx, keys)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get 2fa status: %w", err)
	}
	return val, nil
}
func (us *UserHelper) Clear2FaStatus(ctx context.Context, sessionID, key string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s:%s", TwoFaStatusPrefix, sessionID, key)
	return us.client.DeleteCache(subCtx, keys)
}
func (us *UserHelper) Set2FaCode(ctx context.Context, userID uuid.UUID, secret, qrCode string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s:%s", TwoFaCodePrefix, userID, secret)
	return us.client.SetAny(subCtx, keys, qrCode, 10*time.Minute)
}
func (us *UserHelper) Get2FaCode(ctx context.Context, userID uuid.UUID, secret string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s:%s", TwoFaCodePrefix, userID, secret)
	return us.client.GetString(subCtx, keys)
}
func (us *UserHelper) GetTemp2FASecret(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s", TwoFaCodePrefix, userID)
	val, err := us.client.GetByte(subCtx, keys)
	if err == nil {
		var result map[string]interface{}
		if err = json.Unmarshal([]byte(val), &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Redis data: %v", err)
		}
		return result, nil
	} else if err != redis.Nil {
		return nil, fmt.Errorf("redis error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(val), &payload); err != nil {
		return nil, err
	}
	return payload, nil
}
func (us *UserHelper) Generate2FAQRCode(userID uuid.UUID, email string) (string, string, error) {

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Social-Forge",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	qr, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return "", "", err
	}
	base64QR := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qr)

	// Simpan ke Redis 10 menit
	// err = us.Set2FaCode(subCtx, userID, key.Secret(), base64QR)
	// if err != nil {
	// 	return "", "", err
	// }

	return base64QR, key.Secret(), nil
}
func (us *UserHelper) SetTemp2FASecret(ctx context.Context, userID uuid.UUID, payload map[string]interface{}) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}
	keys := fmt.Sprintf("%s%s", TwoFaCodePrefix, userID)
	return us.client.Setbyte(subCtx, keys, jsonData, 30*time.Minute)
}
func (us *UserHelper) ClearTemp2FASecret(ctx context.Context, userID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("%s%s", TwoFaCodePrefix, userID)
	return us.client.DeleteCache(subCtx, keys)
}
func (us *UserHelper) Verify2FA(ctx context.Context, userID uuid.UUID, otp, secret string) (bool, error) {
	valid, err := totp.ValidateCustom(
		otp,
		secret,
		time.Now(),
		totp.ValidateOpts{
			Period: 30,
			Skew:   1,
			Digits: 6,
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP code: %w", err)
	}

	return valid, nil
}

func (us *UserHelper) SetBlockedAttemptCredential(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	return us.client.SetAny(subCtx, key, val, expiration)
}
func (us *UserHelper) SetExpireAttemptCredential(ctx context.Context, key string, expiration time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	if err := us.client.Expire(subCtx, key, expiration); err != nil {
		return fmt.Errorf("failed to set attempt credential expiration for key %s: %w", key, err)
	}
	return nil
}
func (us *UserHelper) IncrementAndGet(ctx context.Context, key string, expiry time.Duration) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	pipe := us.client.Client().Pipeline()
	incrCmd := pipe.Incr(subCtx, key)
	if expiry > 0 {
		pipe.Expire(subCtx, key, expiry)
	}

	_, err := pipe.Exec(subCtx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment: %w", err)
	}

	return incrCmd.Val(), nil
}
func (us *UserHelper) ResetCounter(ctx context.Context, key string) error {
	return us.client.DeleteCache(ctx, key)
}
func (us *UserHelper) IsBlockedAttempt(ctx context.Context, key string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	val, err := us.client.GetString(subCtx, key)
	if err != nil {
		return false, err
	}

	switch {
	case isPlainNumber(val):
		count, _ := strconv.Atoi(val) // Error sudah dicek di isPlainNumber
		return count >= 3, nil

	case isBinaryData(val):
		count, err := decodeBinaryAttempts([]byte(val))
		if err != nil {
			return false, err
		}
		return count >= 3, nil

	default:
		return false, fmt.Errorf("unrecognized data format for key %s", key)
	}
}
func (us *UserHelper) ShouldBlockCredential(ctx context.Context, key string) int {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	val, err := us.client.GetInt(subCtx, key)
	if err != nil {
		return 0
	}
	return val
}

func (us *UserHelper) ClearSession(ctx context.Context, userID string, key string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keysConfirm := fmt.Sprintf("%s%s", ConfirmPassPerfix, userID)
	if err := us.client.DeleteCache(subCtx, keysConfirm); err != nil {
		return fmt.Errorf("failed to delete confirmation key: %w", err)
	}

	return nil
}
func (us *UserHelper) ClearAllSessionCache(ctx context.Context, userID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	patterns := []string{
		ConfirmPassPerfix + userID,
		TwoFaCodePrefix + userID + "*",
		TwoFaStatusPrefix + "*:" + userID,
		TwoFaStatusPrefix + "*",
		"session:*",
		"revoked_token:*",
	}
	var allErrors []error
	var deletedCount int

	for _, pattern := range patterns {
		count, err := us.deleteKeysByPattern(subCtx, pattern)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to delete cache pattern %s: %w", pattern, err))
			continue
		}
		deletedCount += count
	}

	userSessionCount, err := us.deleteUserSessions(subCtx, userID)
	if err != nil {
		allErrors = append(allErrors, fmt.Errorf("failed to delete user sessions: %w", err))
	} else {
		deletedCount += userSessionCount
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("failed to delete %d cache patterns: %v", deletedCount, allErrors)
	}

	fmt.Printf("✅ Successfully cleared %d cache entries for user %s\n", deletedCount, userID)
	return nil
}
func (us *UserHelper) ClearUserAuthCache(ctx context.Context, userID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	keys := []string{
		ConfirmPassPerfix + userID,
		TwoFaCodePrefix + userID,
	}

	// Tambahkan semua 2FA status keys untuk user ini
	pattern := TwoFaStatusPrefix + "*:" + userID
	statusKeys, err := us.getKeysByPattern(subCtx, pattern)
	if err != nil {
		return fmt.Errorf("failed to get 2FA status keys: %w", err)
	}
	keys = append(keys, statusKeys...)

	if len(keys) > 0 {
		if err := us.client.Client().Del(subCtx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete auth cache: %w", err)
		}
	}

	fmt.Printf("✅ Cleared auth cache for user %s (%d keys)\n", userID, len(keys))
	return nil
}

func (us *UserHelper) deleteKeysByPattern(ctx context.Context, pattern string) (int, error) {
	var deletedCount int
	var cursor uint64
	var lastErr error

	for {
		keys, nextCursor, err := us.client.Client().Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return deletedCount, fmt.Errorf("scan failed for pattern %s: %w", pattern, err)
		}

		if len(keys) > 0 {
			count, err := us.client.Client().Del(ctx, keys...).Result()
			if err != nil {
				lastErr = fmt.Errorf("failed to delete keys: %w", err)
			} else {
				deletedCount += int(count)
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return deletedCount, lastErr
}
func (us *UserHelper) deleteUserSessions(ctx context.Context, userID string) (int, error) {
	var deletedCount int
	var cursor uint64
	pattern := "session:*"

	for {
		keys, nextCursor, err := us.client.Client().Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return deletedCount, fmt.Errorf("scan sessions failed: %w", err)
		}

		for _, key := range keys {
			data, err := us.client.GetByte(ctx, key)
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				continue
			}

			var sessionData struct {
				UserID string `json:"user_id"`
			}
			if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
				continue // Skip invalid session data
			}

			if sessionData.UserID == userID {
				if err := us.client.DeleteCache(ctx, key); err == nil {
					deletedCount++

					sessionID := strings.TrimPrefix(key, "session:")
					revokedKey := "revoked_token:" + sessionID
					us.client.DeleteCache(ctx, revokedKey)
				}
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return deletedCount, nil
}
func (us *UserHelper) getKeysByPattern(ctx context.Context, pattern string) ([]string, error) {
	var allKeys []string
	var cursor uint64

	for {
		keys, nextCursor, err := us.client.Client().Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		allKeys = append(allKeys, keys...)

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return allKeys, nil
}

func isPlainNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func isBinaryData(s string) bool {
	return len(s) > 0 && s[0] == '\x83' // Magic byte MessagePack
}
func decodeBinaryAttempts(data []byte) (int, error) {
	var result struct {
		CurrHits int `msgpack:"currHits"`
	}
	if err := msgpack.Unmarshal(data, &result); err != nil {
		return 0, err
	}
	return result.CurrHits, nil
}
