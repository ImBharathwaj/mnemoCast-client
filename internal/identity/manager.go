package identity

import (
	"encoding/json"
	"fmt"
	"mnemoCast-client/internal/models"
	"os"
	"path/filepath"
	"time"
)

// Manager handles screen identity operations
type Manager struct {
	configDir string
}

// NewManager creates a new identity manager
func NewManager(configDir string) *Manager {
	return &Manager{
		configDir: configDir,
	}
}

// GetOrCreateIdentity loads existing identity
// Note: Identity should be created from server response after connection
// This method no longer auto-generates IDs - screen ID comes from server
func (m *Manager) GetOrCreateIdentity() (*models.ScreenIdentity, error) {
	// Try to load existing identity
	identity, err := m.LoadIdentity()
	if err == nil && identity != nil {
		return identity, nil
	}

	// Identity not found - should be created from server connection
	return nil, fmt.Errorf("identity not found - connect to server first to load screen identity")
}

// LoadIdentity loads the screen identity from file
func (m *Manager) LoadIdentity() (*models.ScreenIdentity, error) {
	identityFile := filepath.Join(m.configDir, "identity.json")
	
	data, err := os.ReadFile(identityFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("identity file not found")
		}
		return nil, fmt.Errorf("failed to read identity file: %w", err)
	}

	var identity models.ScreenIdentity
	if err := json.Unmarshal(data, &identity); err != nil {
		return nil, fmt.Errorf("failed to parse identity: %w", err)
	}

	// Migrate old format if needed (handle nested Location structure)
	var oldFormat struct {
		ID            string    `json:"id"`
		Name          string    `json:"name"`
		Location      struct {
			City      string `json:"city"`
			Area      string `json:"area"`
			VenueType string `json:"venueType"`
			Address   string `json:"address,omitempty"`
		} `json:"location"`
		Classification int       `json:"classification"`
		CreatedAt      time.Time `json:"createdAt"`
		LastSeen       time.Time `json:"lastSeen"`
	}
	
	// Try to parse as old format
	if err := json.Unmarshal(data, &oldFormat); err == nil && oldFormat.Location.City != "" {
		// Migrate from old format
		identity.Country = "Unknown"
		identity.City = oldFormat.Location.City
		identity.Area = oldFormat.Location.Area
		identity.VenueType = oldFormat.Location.VenueType
		if identity.Timezone == "" {
			identity.Timezone = "UTC"
		}
		if identity.Width == 0 {
			identity.Width = 1920
		}
		if identity.Height == 0 {
			identity.Height = 1080
		}
		// Save migrated identity
		m.SaveIdentity(&identity)
	}

	// Validate identity
	if err := identity.Validate(); err != nil {
		return nil, fmt.Errorf("invalid identity: %w", err)
	}

	return &identity, nil
}

// SaveIdentity saves the screen identity to file
func (m *Manager) SaveIdentity(identity *models.ScreenIdentity) error {
	// Ensure config directory exists
	if err := os.MkdirAll(m.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	identityFile := filepath.Join(m.configDir, "identity.json")
	
	data, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal identity: %w", err)
	}

	// Write with restricted permissions
	if err := os.WriteFile(identityFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write identity file: %w", err)
	}

	return nil
}

// CreateIdentityFromServer creates identity from server response
// This is called after successful connection to server
func (m *Manager) CreateIdentityFromServer(screen *models.Screen) (*models.ScreenIdentity, error) {
	now := time.Now()
	identity := &models.ScreenIdentity{
		ID:            screen.ID,
		Name:          screen.Name,
		Country:       screen.Country,
		City:          screen.City,
		Area:          screen.Area,
		VenueType:     screen.VenueType,
		Timezone:      screen.Timezone,
		Width:         screen.Width,
		Height:        screen.Height,
		IsAudible:     screen.IsAudible,
		IsOnline:      screen.IsOnline,
		Classification: screen.Classification,
		CreatedAt:     screen.CreatedAt,
		UpdatedAt:     now,
	}

	// Set LastSeen if available
	if screen.LastSeen != nil {
		identity.LastSeen = screen.LastSeen
	}

	// Save identity
	if err := m.SaveIdentity(identity); err != nil {
		return nil, fmt.Errorf("failed to save identity: %w", err)
	}

	return identity, nil
}

// UpdateIdentityFromServer updates existing identity from server response
func (m *Manager) UpdateIdentityFromServer(screen *models.Screen) error {
	identity, err := m.LoadIdentity()
	if err != nil {
		// If identity doesn't exist, create it
		_, err := m.CreateIdentityFromServer(screen)
		return err
	}

	// Update fields from server response
	identity.Name = screen.Name
	identity.Country = screen.Country
	identity.City = screen.City
	identity.Area = screen.Area
	identity.VenueType = screen.VenueType
	identity.Timezone = screen.Timezone
	identity.Width = screen.Width
	identity.Height = screen.Height
	identity.IsAudible = screen.IsAudible
	identity.IsOnline = screen.IsOnline
	identity.Classification = screen.Classification
	identity.UpdatedAt = time.Now()
	if screen.LastSeen != nil {
		identity.LastSeen = screen.LastSeen
	}

	return m.SaveIdentity(identity)
}

// UpdateLastSeen updates the last seen timestamp
func (m *Manager) UpdateLastSeen(identity *models.ScreenIdentity) error {
	now := time.Now()
	identity.LastSeen = &now
	identity.UpdatedAt = time.Now()
	return m.SaveIdentity(identity)
}

