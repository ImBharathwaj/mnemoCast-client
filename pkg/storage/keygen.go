package storage

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateEncryptionKey generates a random 32-byte key for encryption
func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// GenerateEncryptionKeyBase64 generates a random key and returns it as base64
func GenerateEncryptionKeyBase64() (string, error) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

