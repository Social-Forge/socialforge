package helpers

import (
	"errors"
	"strings"
)

var (
	ErrInvalidKeyLength           = errors.New("key must be 32 bytes for AES-256")
	ErrDecryptionFailed           = errors.New("decryption failed")
	ErrInvalidEncryptionKeyLength = errors.New("encryption key must be 32 bytes (256-bit) after hex decoding")
	ErrFailedToGenerateNonce      = errors.New("failed to generate nonce")
	ErrCiphertextTooShort         = errors.New("ciphertext too short for decryption")
)

type NoRetryError struct {
	error
}
type ValidationErrors struct {
	Errors map[string]string
}

func (e NoRetryError) Error() string {
	return e.error.Error()
}
func (e NoRetryError) Unwrap() error {
	return e.error
}
func NewNoRetryError(msg string) error {
	return NoRetryError{error: errors.New(msg)}
}
func (v ValidationErrors) Error() string {
	var sb strings.Builder
	i := 0
	for field, errMsg := range v.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(field)
		sb.WriteString(": ")
		sb.WriteString(errMsg)
		i++
	}
	return sb.String()
}
