package credentials

import (
	"encoding/json"
	"fmt"
	"mnemoCast-client/internal/models"
	"mnemoCast-client/pkg/storage"
	"os"
	"path/filepath"
)

// Storage handles secure credential storage
type Storage struct {
	configDir string
	keyFile   string
	credsFile string
}

// NewStorage creates a new credential storage
func NewStorage(configDir string) *Storage {
	return &Storage{
		configDir: configDir,
		keyFile:   filepath.Join(configDir, ".encryption_key"),
		credsFile: filepath.Join(configDir, "credentials.json.enc"),
	}
}

// getOrCreateKey gets the encryption key or creates a new one
func (s *Storage) getOrCreateKey() ([]byte, error) {
	// Try to load existing key
	if data, err := os.ReadFile(s.keyFile); err == nil {
		// Key exists, return it
		return data, nil
	}

	// Generate new key
	key, err := storage.GenerateEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Save key with restricted permissions
	if err := os.WriteFile(s.keyFile, key, 0600); err != nil {
		return nil, fmt.Errorf("failed to save encryption key: %w", err)
	}

	return key, nil
}

// Save saves credentials securely
func (s *Storage) Save(creds *models.Credentials) error {
	// Ensure config directory exists
	if err := os.MkdirAll(s.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Get or create encryption key
	key, err := s.getOrCreateKey()
	if err != nil {
		return err
	}

	// Marshal credentials to JSON
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Encrypt and encode to base64 for storage
	encryptedBase64, err := storage.EncryptToBase64(data, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt credentials: %w", err)
	}

	// Save encrypted credentials
	if err := os.WriteFile(s.credsFile, []byte(encryptedBase64), 0600); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

// Load loads credentials from secure storage
func (s *Storage) Load() (*models.Credentials, error) {
	// Check if credentials file exists
	if _, err := os.Stat(s.credsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("credentials file not found")
	}

	// Get encryption key
	key, err := s.getOrCreateKey()
	if err != nil {
		return nil, err
	}

	// Read encrypted credentials
	encryptedBase64, err := os.ReadFile(s.credsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Decrypt credentials
	data, err := storage.DecryptFromBase64(string(encryptedBase64), key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credentials: %w", err)
	}

	// Unmarshal credentials
	var creds models.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}

// Exists checks if credentials file exists
func (s *Storage) Exists() bool {
	_, err := os.Stat(s.credsFile)
	return err == nil
}

// Delete removes credentials file (use with caution)
func (s *Storage) Delete() error {
	if err := os.Remove(s.credsFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete credentials: %w", err)
	}
	return nil
}

