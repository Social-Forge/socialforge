package utils

import (
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(jwtSecret string, sessionData *entity.RedisSessionData, expiry time.Duration) (string, error) {
	now := time.Now().UTC()
	if sessionData == nil {
		return "", fmt.Errorf("session data cannot be nil")
	}
	if sessionData.UserID == uuid.Nil {
		return "", fmt.Errorf("user ID cannot be empty")
	}
	if sessionData.SessionID == "" {
		return "", fmt.Errorf("session ID cannot be empty")
	}

	claims := dto.JWTClaims{
		UserID:       sessionData.UserID.String(),
		Email:        sessionData.Email,
		TenantID:     sessionData.TenantID.String(),
		UserTenantID: sessionData.UserTenantID.String(),
		RoleID:       sessionData.RoleID.String(),
		SessionID:    sessionData.SessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "social-forge",
			Subject:   sessionData.UserID.String(),
			ID:        sessionData.SessionID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
func VerifyJWT(tokenString string, jwtSecret string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token string cannot be empty")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify jwt token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsedToken, nil
}
func GenerateCSRFToken(userId, email string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     time.Now().Add(time.Second * 60).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
func VerifyCSRFToken(tokenString string, jwtSecret string) (*jwt.Token, error) {

	csrfToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected request method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return csrfToken, nil
}
