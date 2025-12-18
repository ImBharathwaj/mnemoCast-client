package config

import (
	"encoding/json"
	"fmt"
	"mnemoCast-client/internal/models"
	"os"
	"path/filepath"
)

// Loader handles configuration loading and saving
type Loader struct {
	configDir string
	configFile string
}

// NewLoader creates a new configuration loader
func NewLoader(configDir string) *Loader {
	return &Loader{
		configDir:  configDir,
		configFile: filepath.Join(configDir, "config.json"),
	}
}

// Load loads the configuration from file
func (l *Loader) Load() (*models.ScreenConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(l.configFile); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return l.CreateDefault()
	}

	data, err := os.ReadFile(l.configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.ScreenConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults for missing fields
	needsSave := false
	if config.AdFetchInterval == 0 {
		config.AdFetchInterval = 60 // Default: fetch ads every 60 seconds (1 minute)
		needsSave = true
	}
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 30 // Default: heartbeat every 30 seconds
		needsSave = true
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = 3 // Default: 3 retry attempts
		needsSave = true
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 5 // Default: 5 seconds between retries
		needsSave = true
	}

	// Save updated config if we added defaults (to persist new fields)
	if needsSave {
		if err := l.Save(&config); err != nil {
			// Log but don't fail - config is still valid
			fmt.Printf("Warning: Failed to save updated config: %v\n", err)
		}
	}

	return &config, nil
}

// Save saves the configuration to file
func (l *Loader) Save(config *models.ScreenConfig) error {
	// Ensure config directory exists
	if err := os.MkdirAll(l.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with restricted permissions
	if err := os.WriteFile(l.configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// CreateDefault creates and saves a default configuration
func (l *Loader) CreateDefault() (*models.ScreenConfig, error) {
	config := models.DefaultConfig()
	
	// Set default ad server URL
	config.AdServerURL = "http://10.42.0.1:8080"
	
	if err := l.Save(config); err != nil {
		return nil, fmt.Errorf("failed to save default config: %w", err)
	}

	return config, nil
}

// GetConfigDir returns the configuration directory path
func (l *Loader) GetConfigDir() string {
	return l.configDir
}

