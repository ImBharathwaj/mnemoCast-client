package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt encrypts data using AES-256-GCM
func Encrypt(data []byte, key []byte) ([]byte, error) {
	// Create hash of key to ensure 32-byte key
	hash := sha256.Sum256(key)
	key32 := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key32)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	// Create hash of key to ensure 32-byte key
	hash := sha256.Sum256(key)
	key32 := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key32)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptToBase64 encrypts data and returns base64 encoded string
func EncryptToBase64(data []byte, key []byte) (string, error) {
	encrypted, err := Encrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptFromBase64 decrypts base64 encoded encrypted data
func DecryptFromBase64(encryptedBase64 string, key []byte) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return Decrypt(encrypted, key)
}

// DeriveKeyFromPassword derives an encryption key from a password
func DeriveKeyFromPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

