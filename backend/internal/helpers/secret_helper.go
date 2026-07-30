package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"go.uber.org/zap"
)

type SecretHelper struct {
	encryptionKey []byte
	logger        *zap.Logger
}

type SecretService interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

func NewSecretHelper(
	encryptionKeyHex string,
	logger *zap.Logger,
) (*SecretHelper, error) {
	keyBytes, err := hex.DecodeString(encryptionKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid hex encryption key provided: %w", err)
	}
	if len(keyBytes) != 32 {
		return nil, ErrInvalidEncryptionKeyLength
	}

	return &SecretHelper{
		encryptionKey: keyBytes,
		logger:        logger,
	}, nil
}
func (h *SecretHelper) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}
func (h *SecretHelper) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(h.encryptionKey)
	if err != nil {
		h.logger.Error("Failed to create AES cipher", zap.Error(err))
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		h.logger.Error("Failed to create GCM cipher", zap.Error(err))
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		h.logger.Error("Failed to generate nonce", zap.Error(err))
		return nil, ErrFailedToGenerateNonce
	}
	// Seal appends the ciphertext and the authentication tag to the nonce.
	// The final output format is: nonce || ciphertext || tag
	ciphertextWithNonce := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertextWithNonce, nil
}
func (h *SecretHelper) Decrypt(ciphertextWithNonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(h.encryptionKey)
	if err != nil {
		h.logger.Error("Failed to create AES cipher for decryption", zap.Error(err))
		return nil, fmt.Errorf("failed to create AES cipher for decryption: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		h.logger.Error("Failed to create GCM cipher for decryption", zap.Error(err))
		return nil, fmt.Errorf("failed to create GCM cipher for decryption: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextWithNonce) < nonceSize {
		h.logger.Error("Ciphertext too short for decryption", zap.Int("length", len(ciphertextWithNonce)), zap.Int("nonce_size", nonceSize))
		return nil, ErrCiphertextTooShort
	}

	nonce, ciphertext := ciphertextWithNonce[:nonceSize], ciphertextWithNonce[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		h.logger.Error("Failed to decrypt data", zap.Error(err))
		// The underlying error might be crypto/cipher: message authentication failed
		// which indicates tampering or wrong key.
		return nil, ErrDecryptionFailed
	}
	return plaintext, nil
}
