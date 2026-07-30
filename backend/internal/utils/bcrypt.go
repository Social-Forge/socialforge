package utils

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}
	return string(hash), nil
}
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	// Contains uppercase, lowercase, digit, and special character
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false
	}

	// Basic check for common weak passwords (expand this list)
	weakPasswords := []string{"password", "123456", "qwerty", "12345678", "12345", "1234"}
	for _, weak := range weakPasswords {
		if strings.Contains(strings.ToLower(password), weak) {
			return false
		}
	}

	return true
}
