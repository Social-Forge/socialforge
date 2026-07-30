package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

func GenerateEmailToken() (string, time.Time) {
	return uuid.New().String(), time.Now().Add(24 * time.Hour)
}
func GenerateLicense() string {
	uuidStr := uuid.New().String()
	timestamp := time.Now().UnixNano()
	uuidWithTImestamp := fmt.Sprintf("sag-%s%x", uuidStr[:24], timestamp)
	return uuidWithTImestamp
}
func GenerateTaskNumber(siteID uuid.UUID) string {
	uuidStr := siteID.String()
	timestamp := time.Now().UnixNano()
	uuidWithTImestamp := fmt.Sprintf("task-%s%x", uuidStr[:24], timestamp)
	return uuidWithTImestamp
}
func GenerateBatchNumber() string {
	uuidStr := uuid.New().String()
	timestamp := time.Now().UnixNano()
	uuidWithTImestamp := fmt.Sprintf("batch-%s%x", uuidStr[:24], timestamp)
	return uuidWithTImestamp
}
func ParseLicense(license string) (string, time.Time, error) {
	parts := strings.Split(license, "|")
	if len(parts) != 2 {
		return "", time.Time{}, fmt.Errorf("invalid license format")
	}

	expiry, err := time.Parse("2006-01-02T15:04:05", parts[1])
	if err != nil {
		return "", time.Time{}, err
	}

	return parts[0], expiry, nil
}
func GenerateRandomString(length int) (string, error) {
	if length < 16 { // Minimal 16 karakter untuk memastikan ada ruang untuk random data
		return "", errors.New("length must be at least 16")
	}
	// 1. Ambil timestamp nano (16 karakter hex)
	timestamp := time.Now().UnixNano()
	timestampHex := []byte(fmt.Sprintf("%016x", timestamp))
	// 2. Generate random bytes (length - 16)
	randomBytes := make([]byte, (length-16)*6/8) // Base64 overhead
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	// 3. Gabungkan dan encode
	combined := append(timestampHex, randomBytes...)
	encoded := base64.URLEncoding.EncodeToString(combined)
	// 4. Pastikan panjang sesuai
	if len(encoded) > length {
		encoded = encoded[:length]
	}
	return encoded, nil
}
func GenerateAPIKey() string {
	// 32 bytes random (43 chars ketika di-encode base64)
	b := make([]byte, 32)
	rand.Read(b)
	return "ak_" + base64.URLEncoding.EncodeToString(b)
}
func GenerateVerificationCode() (string, error) {
	// 1. Get last 8 digits of nanosecond timestamp
	timestamp := time.Now().UnixNano() % 1e8 // Last 8 digits
	// 2. Generate secure random 4-digit number
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return "", err
	}
	randomNum := n.Int64() + 1000 // Ensure 4 digits (1000-9999)

	return fmt.Sprintf("VC-%08d-%04d", timestamp, randomNum), nil
}
func GenerateRefreshToken() (string, time.Time) {
	token := uuid.New().String()
	expires := time.Now().Add(7 * 24 * time.Hour) // 7 hari
	return token, expires
}
func GenerateRandomToken() string {
	uuidStr := uuid.New().String()
	timestamp := time.Now().UnixNano()
	uuidWithTImestamp := fmt.Sprintf("%s%x", uuidStr[:24], timestamp)
	return uuidWithTImestamp
}
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err) // handle error properly in production
	}
	return base64.URLEncoding.EncodeToString(b)
}
func GenerateRandomPassword() (string, error) {
	strRes, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to generate password : %w", err)
	}
	return strRes, nil
}
