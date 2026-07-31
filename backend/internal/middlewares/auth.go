package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/metrics"
	redisclient "github/socialforge/internal/infra/redis-client"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrMissingToken   = "Unauthorized - Missing token"
	ErrMissMatchToken = "Unauthorized - Token mismatch"
	ErrInvalidToken   = "Unauthorized - Invalid token format"
	ErrExpiredToken   = "Unauthorized - Invalid or expired token"
	ErrInvalidAuth    = "Unauthorized - Invalid authorization header"
)

type AuthMiddleware struct {
	notifier    config.Notifier
	ctxinject   *ContextMiddleware
	redisClient *redisclient.RedisClient
	tokenHelper *helpers.TokenHelper
	logger      *zap.Logger
	jwtSecret   string
}

func NewAuthMiddleware(notifier config.Notifier, ctxinject *ContextMiddleware, redisClient *redisclient.RedisClient, tokenHelper *helpers.TokenHelper, logger *zap.Logger, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		notifier:    notifier,
		ctxinject:   ctxinject,
		redisClient: redisClient,
		tokenHelper: tokenHelper,
		logger:      logger,
		jwtSecret:   jwtSecret,
	}
}

// RequireRoles guards a route so only callers holding one of the given role
// names may pass. Must run AFTER JWTAuth (which populates role_name in Locals).
func (m *AuthMiddleware) RequireRoles(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c fiber.Ctx) error {
		names, ok := c.Locals("role_name").([]string)
		if !ok {
			return helpers.Respond(c, fiber.StatusForbidden, "Insufficient permissions", nil)
		}
		for _, n := range names {
			if _, found := allowed[n]; found {
				return c.Next()
			}
		}
		return helpers.Respond(c, fiber.StatusForbidden, "Insufficient permissions", nil)
	}
}

func (m *AuthMiddleware) JWTAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := m.ctxinject.From(c)

		tokenStr, err := m.extractTokenFromHeader(c)
		if err != nil {
			return helpers.Respond(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		m.logger.Debug("🔐 [STEP 1] Extracted token from header",
			zap.String("token_prefix", tokenStr[:20]+"..."),
			zap.String("path", c.Path()))

		jwtClaims, sessionID, err := m.validateJWTAndExtractSession(tokenStr)
		if err != nil {
			m.logger.Error("❌ [STEP 1 FAILED] JWT validation failed",
				zap.Error(err),
				zap.String("token_prefix", tokenStr[:20]+"..."))
			return m.handleTokenError(c, err)
		}
		m.logger.Debug("✅ [STEP 1 SUCCESS] JWT validated",
			zap.String("session_id", sessionID),
			zap.String("user_id", jwtClaims["sub"].(string)))

		redisSessionData, err := m.validateTokenSession(ctx, sessionID)
		if err != nil {
			m.logger.Error("❌ [STEP 2 FAILED] Redis session validation failed",
				zap.Error(err),
				zap.String("session_id", sessionID),
				zap.String("redis_key", "session:"+sessionID))
			return helpers.Respond(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		m.logger.Debug("✅ [STEP 2 SUCCESS] Redis session found",
			zap.String("user_id", redisSessionData.UserID.String()),
			zap.String("email", redisSessionData.Email))

		if err := m.crossValidateClaims(jwtClaims, redisSessionData); err != nil {
			m.logger.Warn("❌ [STEP 3 FAILED] Token claim mismatch",
				zap.String("user_id", redisSessionData.UserID.String()),
				zap.Error(err))
			return helpers.Respond(c, fiber.StatusUnauthorized, ErrInvalidToken, nil)
		}

		m.logger.Debug("✅ [STEP 3 SUCCESS] Claims cross-validated")

		m.setContextLocals(c, redisSessionData)

		m.logger.Debug("✅ [FINAL] JWT authentication successful",
			zap.String("user_id", redisSessionData.UserID.String()),
			zap.String("session_id", sessionID),
			zap.String("path", c.Path()))

		return c.Next()
	}
}
func (m *AuthMiddleware) extractTokenFromHeader(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New(ErrMissingToken)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New(ErrInvalidAuth)
	}

	if parts[1] == "" {
		return "", errors.New(ErrMissingToken)
	}

	return parts[1], nil
}
func (m *AuthMiddleware) validateJWTAndExtractSession(tokenStr string) (jwt.MapClaims, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, "", err
	}

	if !token.Valid {
		return nil, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", errors.New("invalid token claims")
	}

	if exp, okExp := claims["exp"].(float64); okExp {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return nil, "", jwt.ErrTokenExpired
		}
	} else {
		return nil, "", errors.New("missing expiration claim")
	}

	if _, okSub := claims["sub"].(string); !okSub {
		return nil, "", errors.New("missing subject claim")
	}

	sessionID, okSid := claims["sid"].(string)
	if !okSid || sessionID == "" {
		return nil, "", errors.New("missing session ID in token")
	}

	return claims, sessionID, nil
}
func (m *AuthMiddleware) validateTokenSession(ctx context.Context, sessionID string) (*entity.RedisSessionData, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	m.logger.Debug("🔍 [validateTokenSession] Looking up session",
		zap.String("session_id", sessionID))

	redisSessionData, err := m.tokenHelper.GetSessionTokenMetadata(subCtx, sessionID)
	if err != nil {
		m.logger.Error("❌ [validateTokenSession] GetSessionTokenMetadata failed",
			zap.Error(err),
			zap.String("session_id", sessionID))
		return nil, errors.New(ErrInvalidToken)
	}

	m.logger.Debug("🔍 [validateTokenSession] Redis data retrieved",
		zap.String("user_id", redisSessionData.UserID.String()),
		zap.String("email", redisSessionData.Email))

	if err := m.validateTokenMetadata(redisSessionData); err != nil {
		m.logger.Error("❌ [validateTokenSession] validateTokenMetadata failed",
			zap.Error(err),
			zap.String("user_id", redisSessionData.UserID.String()))
		return nil, err
	}

	return redisSessionData, nil
}

func (m *AuthMiddleware) handleTokenError(c fiber.Ctx, err error) error {
	errorMsg := ErrInvalidToken
	if errors.Is(err, jwt.ErrTokenExpired) {
		errorMsg = ErrExpiredToken
	}

	// Log metrics untuk monitoring
	metrics.GetAppMetrics().JWTErrorTotal.WithLabelValues("validation_failed").Inc()

	return helpers.Respond(c, fiber.StatusUnauthorized, errorMsg, nil)
}
func (m *AuthMiddleware) crossValidateClaims(jwtClaims jwt.MapClaims, redisData *entity.RedisSessionData) error {
	jwtUserIDStr, ok := jwtClaims["sub"].(string)
	if !ok {
		return errors.New("missing user ID in JWT claims")
	}

	jwtUserID, err := uuid.Parse(jwtUserIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID format in JWT: %w", err)
	}

	if jwtUserID != redisData.UserID {
		return fmt.Errorf("user ID mismatch: JWT=%s, Redis=%s", jwtUserID, redisData.UserID)
	}

	// ✅ PERUBAHAN: Validasi session ID wajib
	if jwtSessionID, ok := jwtClaims["sid"].(string); ok && jwtSessionID != "" {
		if jwtSessionID != redisData.SessionID {
			return fmt.Errorf("session ID mismatch: JWT=%s, Redis=%s", jwtSessionID, redisData.SessionID)
		}
	} else {
		return errors.New("missing session ID in JWT claims")
	}

	// ✅ PERUBAHAN: Validasi role ID wajib
	if jwtRoleIDStr, ok := jwtClaims["rid"].(string); ok && jwtRoleIDStr != "" {
		jwtRoleID, err := uuid.Parse(jwtRoleIDStr)
		if err != nil {
			return fmt.Errorf("invalid role ID in JWT: %w", err)
		}
		if jwtRoleID != redisData.RoleID {
			return fmt.Errorf("role ID mismatch: JWT=%s, Redis=%s", jwtRoleID, redisData.RoleID)
		}
	} else {
		return errors.New("missing role ID in JWT claims")
	}

	// ✅ OPTIONAL: Validasi email jika ada
	if jwtEmail, ok := jwtClaims["em"].(string); ok && jwtEmail != "" {
		if jwtEmail != redisData.Email {
			return fmt.Errorf("email mismatch: JWT=%s, Redis=%s", jwtEmail, redisData.Email)
		}
	}

	// ✅ OPTIONAL: Validasi tenant ID jika ada di kedua sisi
	if jwtTenantIDStr, ok := jwtClaims["tid"].(string); ok && jwtTenantIDStr != "" {
		jwtTenantID, err := uuid.Parse(jwtTenantIDStr)
		if err != nil {
			m.logger.Warn("Invalid tenant ID in JWT",
				zap.String("jwt_tenant_id", jwtTenantIDStr),
				zap.Error(err))
			// Tidak return error, hanya log warning
		} else if redisData.TenantID != uuid.Nil && jwtTenantID != redisData.TenantID {
			return fmt.Errorf("tenant ID mismatch: JWT=%s, Redis=%s", jwtTenantID, redisData.TenantID)
		}
	}

	// ✅ OPTIONAL: Validasi user tenant ID jika ada di kedua sisi
	if jwtUserTenantIDStr, ok := jwtClaims["utid"].(string); ok && jwtUserTenantIDStr != "" {
		jwtUserTenantID, err := uuid.Parse(jwtUserTenantIDStr)
		if err != nil {
			m.logger.Warn("Invalid user tenant ID in JWT",
				zap.String("jwt_user_tenant_id", jwtUserTenantIDStr),
				zap.Error(err))
			// Tidak return error, hanya log warning
		} else if redisData.UserTenantID != uuid.Nil && jwtUserTenantID != redisData.UserTenantID {
			return fmt.Errorf("user tenant ID mismatch: JWT=%s, Redis=%s", jwtUserTenantID, redisData.UserTenantID)
		}
	}

	return nil
}
func (m *AuthMiddleware) validateTokenMetadata(metadata *entity.RedisSessionData) error {
	if metadata.UserID == uuid.Nil {
		return fmt.Errorf("invalid user ID in session")
	}
	if metadata.RoleID == uuid.Nil {
		return fmt.Errorf("invalid role ID in session")
	}
	if metadata.SessionID == "" {
		return fmt.Errorf("invalid session ID")
	}
	if metadata.TenantID == uuid.Nil {
		m.logger.Warn("Tenant ID is nil in session",
			zap.String("user_id", metadata.UserID.String()),
			zap.String("session_id", metadata.SessionID))
	}
	if metadata.UserTenantID == uuid.Nil {
		m.logger.Warn("User tenant ID is nil in session",
			zap.String("user_id", metadata.UserID.String()),
			zap.String("session_id", metadata.SessionID))
	}

	return nil
}
func (m *AuthMiddleware) setContextLocals(c fiber.Ctx, metadata *entity.RedisSessionData) {
	c.Locals("user_id", metadata.UserID.String())
	c.Locals("tenant_id", metadata.TenantID.String())
	c.Locals("user_tenant_id", metadata.UserTenantID.String())
	c.Locals("role_id", metadata.RoleID.String()) // Pakai RoleID, bukan Role.ID
	c.Locals("email", metadata.Email)
	c.Locals("session_id", metadata.SessionID)

	if len(metadata.RoleName) > 0 {
		c.Locals("role_name", metadata.RoleName)
	}
}
func (am *AuthMiddleware) LogUnauthorized(c fiber.Ctx, subject string, requestID string) {
	metrics.GetAppMetrics().JWTErrorTotal.WithLabelValues("invalid_signature").Inc()
	am.notifier.SendAlert(config.AlertRequest{
		Subject: subject,
		Message: subject,
		Metadata: map[string]interface{}{
			"request_id": requestID,
			"timestamp":  time.Now(),
			"user_agent": c.Get("User-Agent"),
			"ip":         c.Locals("real_ip").(string),
		},
	})
}
