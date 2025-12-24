package player

import (
	"fmt"
	"io"
	"log"
	"mnemoCast-client/internal/ads"
	"mnemoCast-client/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Downloader handles downloading and caching ad media files
type Downloader struct {
	storage    *ads.Storage
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
}

// NewDownloader creates a new downloader instance
func NewDownloader(storage *ads.Storage, maxRetries int, retryDelaySeconds int) *Downloader {
	return &Downloader{
		storage: storage,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: maxRetries,
		retryDelay: time.Duration(retryDelaySeconds) * time.Second,
	}
}

// DownloadAdMedia downloads the media file for an ad
// Returns the local file path if successful
// Supports both HTTP URLs and file:// URLs for local testing
func (d *Downloader) DownloadAdMedia(ad *models.Ad) (string, error) {
	// Check if already cached
	if localPath, exists := d.GetLocalPath(ad); exists {
		log.Printf("[%s] [DOWNLOAD] Media already cached: %s", time.Now().Format("15:04:05.000"), localPath)
		return localPath, nil
	}
	
	// Handle file:// URLs (for local testing)
	if strings.HasPrefix(ad.ContentURL, "file://") {
		localPath := strings.TrimPrefix(ad.ContentURL, "file://")
		if _, err := os.Stat(localPath); err == nil {
			log.Printf("[%s] [DOWNLOAD] Using local file: %s", time.Now().Format("15:04:05.000"), localPath)
			return localPath, nil
		}
		return "", fmt.Errorf("local file not found: %s", localPath)
	}
	
	// Ensure ad media directory exists
	if err := d.storage.EnsureAdMediaDir(ad.ID); err != nil {
		return "", fmt.Errorf("failed to create media directory: %w", err)
	}
	
	// Determine file extension from URL or ad type
	ext := d.getFileExtension(ad.ContentURL, ad.Type)
	fileName := fmt.Sprintf("%s%s", ad.ID, ext)
	localPath := d.storage.GetAdMediaPath(ad.ID, fileName)
	
	// Download file
	log.Printf("[%s] [DOWNLOAD] Downloading media: %s -> %s", 
		time.Now().Format("15:04:05.000"), ad.ContentURL, localPath)
	
	var lastErr error
	for attempt := 0; attempt <= d.maxRetries; attempt++ {
		if attempt > 0 {
			delay := d.retryDelay * time.Duration(attempt)
			log.Printf("[%s] [DOWNLOAD] Retrying download (attempt %d/%d) after %v...", 
				time.Now().Format("15:04:05.000"), attempt, d.maxRetries, delay)
			time.Sleep(delay)
		}
		
		err := d.DownloadFile(ad.ContentURL, localPath)
		if err == nil {
			// Verify file was downloaded successfully
			if info, err := os.Stat(localPath); err == nil && info.Size() > 0 {
				log.Printf("[%s] [DOWNLOAD] Media downloaded successfully: %s (%d bytes)", 
					time.Now().Format("15:04:05.000"), localPath, info.Size())
				return localPath, nil
			}
			err = fmt.Errorf("downloaded file is empty or invalid")
		}
		
		lastErr = err
		log.Printf("[%s] [DOWNLOAD] Download attempt %d/%d failed: %v", 
			time.Now().Format("15:04:05.000"), attempt+1, d.maxRetries+1, err)
	}
	
	return "", fmt.Errorf("failed to download media after %d attempts: %w", d.maxRetries+1, lastErr)
}

// GetLocalPath returns the local path for an ad's media file if it exists
func (d *Downloader) GetLocalPath(ad *models.Ad) (string, bool) {
	ext := d.getFileExtension(ad.ContentURL, ad.Type)
	fileName := fmt.Sprintf("%s%s", ad.ID, ext)
	localPath := d.storage.GetAdMediaPath(ad.ID, fileName)
	
	if info, err := os.Stat(localPath); err == nil && info.Size() > 0 {
		return localPath, true
	}
	
	return "", false
}

// IsCached checks if an ad's media is already cached
func (d *Downloader) IsCached(ad *models.Ad) bool {
	_, exists := d.GetLocalPath(ad)
	return exists
}

// DownloadFile downloads a file from a URL to a local path
func (d *Downloader) DownloadFile(url, destPath string) error {
	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set user agent
	req.Header.Set("User-Agent", "MnemoCast-Client/1.0")
	
	// Execute request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	// Create destination file
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()
	
	// Copy response body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// getFileExtension determines the file extension from URL or ad type
func (d *Downloader) getFileExtension(url, adType string) string {
	// Try to extract extension from URL
	if ext := filepath.Ext(url); ext != "" {
		// Remove query parameters
		ext = strings.Split(ext, "?")[0]
		if ext != "" {
			return ext
		}
	}
	
	// Fall back to ad type
	switch strings.ToLower(adType) {
	case "image":
		return ".jpg"
	case "video":
		return ".mp4"
	case "html":
		return ".html"
	case "text":
		return ".txt"
	default:
		return ".bin"
	}
}

// CleanupOldMedia removes media files for ads that are no longer in the playlist
func (d *Downloader) CleanupOldMedia(currentAdIDs map[string]bool) error {
	mediaDir := d.storage.GetMediaDir()
	
	// Read all directories in media folder
	entries, err := os.ReadDir(mediaDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Media directory doesn't exist, nothing to clean
		}
		return fmt.Errorf("failed to read media directory: %w", err)
	}
	
	cleaned := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		adID := entry.Name()
		if !currentAdIDs[adID] {
			// Ad no longer exists, remove its media directory
			adMediaDir := filepath.Join(mediaDir, adID)
			if err := os.RemoveAll(adMediaDir); err != nil {
				log.Printf("[%s] [DOWNLOAD] Failed to remove old media directory %s: %v", 
					time.Now().Format("15:04:05.000"), adMediaDir, err)
			} else {
				cleaned++
				log.Printf("[%s] [DOWNLOAD] Cleaned up old media: %s", 
					time.Now().Format("15:04:05.000"), adMediaDir)
			}
		}
	}
	
	if cleaned > 0 {
		log.Printf("[%s] [DOWNLOAD] Cleaned up %d old media directories", 
			time.Now().Format("15:04:05.000"), cleaned)
	}
	
	return nil
}

