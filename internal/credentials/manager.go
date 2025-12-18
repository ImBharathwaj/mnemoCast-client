package credentials

import (
	"fmt"
	"mnemoCast-client/internal/models"
)

// Manager handles credential operations
type Manager struct {
	storage *Storage
}

// NewManager creates a new credential manager
func NewManager(configDir string) *Manager {
	return &Manager{
		storage: NewStorage(configDir),
	}
}

// GetOrCreate loads existing credentials or creates new ones
func (m *Manager) GetOrCreate(screenID string) (*models.Credentials, error) {
	// Try to load existing credentials
	if m.storage.Exists() {
		creds, err := m.storage.Load()
		if err == nil {
			return creds, nil
		}
	}

	// Create new credentials (without passkey - must be set separately)
	creds := &models.Credentials{
		ScreenID: screenID,
	}

	return creds, nil
}

// Load loads credentials from storage
func (m *Manager) Load() (*models.Credentials, error) {
	if !m.storage.Exists() {
		return nil, models.ErrCredentialsNotFound
	}

	creds, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	return creds, nil
}

// Save saves credentials to storage
func (m *Manager) Save(creds *models.Credentials) error {
	// Validate credentials
	if !creds.IsValid() {
		return models.ErrInvalidCredentials
	}

	if err := m.storage.Save(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

// SetCredentials sets both screen ID and passkey in credentials
func (m *Manager) SetCredentials(screenID string, passkey string) error {
	creds := &models.Credentials{
		ScreenID: screenID,
		Passkey:  passkey,
	}

	return m.Save(creds)
}

// GetCredentials gets both screen ID and passkey from credentials
func (m *Manager) GetCredentials() (screenID string, passkey string, err error) {
	creds, err := m.Load()
	if err != nil {
		return "", "", err
	}

	if !creds.HasCredentials() {
		return "", "", models.ErrCredentialsNotFound
	}

	return creds.ScreenID, creds.Passkey, nil
}

// Validate checks if credentials are valid
func (m *Manager) Validate() error {
	creds, err := m.Load()
	if err != nil {
		return err
	}

	if !creds.IsValid() {
		return models.ErrInvalidCredentials
	}

	return nil
}

// Exists checks if credentials exist
func (m *Manager) Exists() bool {
	return m.storage.Exists()
}

// Delete removes credentials (use with caution)
func (m *Manager) Delete() error {
	return m.storage.Delete()
}

