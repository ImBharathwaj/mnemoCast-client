package ads

import (
	"encoding/json"
	"fmt"
	"mnemoCast-client/internal/models"
	"os"
	"path/filepath"
	"time"
)

// Storage handles ad storage in the filesystem
type Storage struct {
	adsDir      string
	adsFile     string
	mediaDir    string
}

// NewStorage creates a new ad storage
func NewStorage(configDir string) *Storage {
	adsDir := filepath.Join(configDir, "ads")
	return &Storage{
		adsDir:   adsDir,
		adsFile:  filepath.Join(adsDir, "current_ads.json"),
		mediaDir: filepath.Join(adsDir, "media"),
	}
}

// SaveAds saves the fetched ads to the filesystem
func (s *Storage) SaveAds(adResponse *models.AdDeliveryResponse) error {
	// Ensure ads directory exists
	if err := os.MkdirAll(s.adsDir, 0755); err != nil {
		return fmt.Errorf("failed to create ads directory: %w", err)
	}

	// Ensure media directory exists
	if err := os.MkdirAll(s.mediaDir, 0755); err != nil {
		return fmt.Errorf("failed to create media directory: %w", err)
	}

	// Create metadata structure with fetch timestamp
	adsMetadata := struct {
		FetchedAt   time.Time                `json:"fetchedAt"`
		PlaylistID  string                   `json:"playlistId,omitempty"`
		UpdatedAt   time.Time                `json:"updatedAt"`
		Ads         []models.Ad              `json:"ads"`
		AdsCount    int                      `json:"adsCount"`
	}{
		FetchedAt:  time.Now(),
		PlaylistID: adResponse.PlaylistID,
		UpdatedAt:  adResponse.UpdatedAt,
		Ads:        adResponse.Ads,
		AdsCount:   len(adResponse.Ads),
	}

	// Save ads metadata to JSON file
	data, err := json.MarshalIndent(adsMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ads: %w", err)
	}

	// Write with restricted permissions
	if err := os.WriteFile(s.adsFile, data, 0600); err != nil {
		return fmt.Errorf("failed to save ads file: %w", err)
	}

	return nil
}

// LoadAds loads ads from the filesystem
func (s *Storage) LoadAds() (*models.AdDeliveryResponse, error) {
	// Check if ads file exists
	if _, err := os.Stat(s.adsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("ads file not found")
	}

	data, err := os.ReadFile(s.adsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ads file: %w", err)
	}

	var adsMetadata struct {
		FetchedAt   time.Time                `json:"fetchedAt"`
		PlaylistID  string                   `json:"playlistId,omitempty"`
		UpdatedAt   time.Time                `json:"updatedAt"`
		Ads         []models.Ad              `json:"ads"`
		AdsCount    int                      `json:"adsCount"`
	}

	if err := json.Unmarshal(data, &adsMetadata); err != nil {
		return nil, fmt.Errorf("failed to parse ads file: %w", err)
	}

	return &models.AdDeliveryResponse{
		PlaylistID: adsMetadata.PlaylistID,
		UpdatedAt:  adsMetadata.UpdatedAt,
		Ads:        adsMetadata.Ads,
	}, nil
}

// GetMediaDir returns the media directory path
func (s *Storage) GetMediaDir() string {
	return s.mediaDir
}

// GetAdsDir returns the ads directory path
func (s *Storage) GetAdsDir() string {
	return s.adsDir
}

// Exists checks if ads file exists
func (s *Storage) Exists() bool {
	_, err := os.Stat(s.adsFile)
	return err == nil
}

// GetAdMediaPath returns the local filesystem path for an ad's media file
func (s *Storage) GetAdMediaPath(adID, fileName string) string {
	return filepath.Join(s.mediaDir, adID, fileName)
}

// EnsureAdMediaDir ensures the directory for a specific ad's media exists
func (s *Storage) EnsureAdMediaDir(adID string) error {
	adMediaDir := filepath.Join(s.mediaDir, adID)
	return os.MkdirAll(adMediaDir, 0755)
}

