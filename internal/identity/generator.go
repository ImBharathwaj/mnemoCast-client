package identity

import (
	"crypto/rand"
	"fmt"
)

// GenerateUUID generates a UUID v4
func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Set version (4) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// GenerateScreenID generates a screen ID with prefix
func GenerateScreenID() (string, error) {
	uuid, err := GenerateUUID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("screen-%s", uuid), nil
}

