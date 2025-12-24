package main

import (
	"fmt"
	"log"
	"mnemoCast-client/internal/ads"
	"mnemoCast-client/internal/models"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/test-ads/main.go <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  create-sample  - Create sample ads for testing")
		fmt.Println("  add-image      - Add an image ad (requires: <ad-id> <image-path> [title] [duration])")
		fmt.Println("  add-video      - Add a video ad (requires: <ad-id> <video-path> [title] [duration])")
		fmt.Println("  add-text       - Add a text ad (requires: <ad-id> <text-content> [title] [duration])")
		fmt.Println("  list           - List current ads")
		fmt.Println("  clear          - Clear all ads")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run cmd/test-ads/main.go create-sample")
		fmt.Println("  go run cmd/test-ads/main.go add-image ad-001 /path/to/image.jpg 'My Ad' 30")
		fmt.Println("  go run cmd/test-ads/main.go add-video ad-002 /path/to/video.mp4 'Video Ad' 60")
		fmt.Println("  go run cmd/test-ads/main.go add-text ad-003 'Hello World' 'Text Ad' 15")
		fmt.Println("  go run cmd/test-ads/main.go list")
		fmt.Println("  go run cmd/test-ads/main.go clear")
		os.Exit(1)
	}

	// Get config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	configDir := filepath.Join(homeDir, ".mnemocast")
	storage := ads.NewStorage(configDir)

	command := os.Args[1]

	switch command {
	case "create-sample":
		createSampleAds(storage)
	case "add-image":
		if len(os.Args) < 4 {
			log.Fatal("Usage: add-image <ad-id> <image-path> [title] [duration]")
		}
		addImageAd(storage, os.Args[2], os.Args[3], getOptionalArg(os.Args, 4, ""), getOptionalIntArg(os.Args, 5, 30))
	case "add-video":
		if len(os.Args) < 4 {
			log.Fatal("Usage: add-video <ad-id> <video-path> [title] [duration]")
		}
		addVideoAd(storage, os.Args[2], os.Args[3], getOptionalArg(os.Args, 4, ""), getOptionalIntArg(os.Args, 5, 60))
	case "add-text":
		if len(os.Args) < 4 {
			log.Fatal("Usage: add-text <ad-id> <text-content> [title] [duration]")
		}
		addTextAd(storage, os.Args[2], os.Args[3], getOptionalArg(os.Args, 4, ""), getOptionalIntArg(os.Args, 5, 15))
	case "list":
		listAds(storage)
	case "clear":
		clearAds(storage)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func createSampleAds(storage *ads.Storage) {
	fmt.Println("Creating sample ads for testing...")
	fmt.Println()

	now := time.Now()
	ads := []models.Ad{
		{
			ID:         "sample-image-001",
			Title:      "Sample Image Ad",
			Type:       "image",
			ContentURL: "file:///tmp/sample-image.jpg", // Placeholder - will be replaced with local path
			Duration:   30,
			Priority:   1,
		},
		{
			ID:         "sample-text-001",
			Title:      "Sample Text Ad",
			Type:       "text",
			ContentURL: "file:///tmp/sample-text.txt",
			Duration:   15,
			Priority:   2,
		},
		{
			ID:         "sample-video-001",
			Title:      "Sample Video Ad",
			Type:       "video",
			ContentURL: "file:///tmp/sample-video.mp4",
			Duration:   45,
			Priority:   1,
		},
	}

	// Create sample media files
	mediaDir := storage.GetMediaDir()
	os.MkdirAll(mediaDir, 0755)

	// Create a simple text file
	textAdDir := filepath.Join(mediaDir, "sample-text-001")
	os.MkdirAll(textAdDir, 0755)
	textFile := filepath.Join(textAdDir, "sample-text-001.txt")
	if err := os.WriteFile(textFile, []byte("This is a sample text advertisement.\n\nWelcome to MnemoCast!\n\nThis ad will display for 15 seconds."), 0644); err == nil {
		ads[1].ContentURL = "file://" + textFile
		fmt.Printf("[OK] Created sample text file: %s\n", textFile)
	}

	response := &models.AdDeliveryResponse{
		Ads:       ads,
		PlaylistID: "sample-playlist",
		UpdatedAt:  now,
	}

	if err := storage.SaveAds(response); err != nil {
		log.Fatalf("Failed to save sample ads: %v", err)
	}

	fmt.Println()
	fmt.Printf("[OK] Created %d sample ads\n", len(ads))
	fmt.Println()
	fmt.Println("Note: For image and video ads, you need to:")
	fmt.Println("  1. Place your media files in the media directory")
	fmt.Println("  2. Update the ContentURL in the ads file")
	fmt.Println()
	fmt.Printf("Media directory: %s\n", mediaDir)
	fmt.Println()
	fmt.Println("To add real media files, use:")
	fmt.Println("  go run cmd/test-ads/main.go add-image <ad-id> <image-path>")
	fmt.Println("  go run cmd/test-ads/main.go add-video <ad-id> <video-path>")
}

func addImageAd(storage *ads.Storage, adID, imagePath, title string, duration int) {
	// Check if image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Image file not found: %s", imagePath)
	}

	// Copy image to media directory
	mediaDir := storage.GetMediaDir()
	adMediaDir := filepath.Join(mediaDir, adID)
	os.MkdirAll(adMediaDir, 0755)

	ext := filepath.Ext(imagePath)
	destFile := filepath.Join(adMediaDir, adID+ext)

	// Read and copy file
	data, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}

	if err := os.WriteFile(destFile, data, 0644); err != nil {
		log.Fatalf("Failed to copy image file: %v", err)
	}

	// Load existing ads
	existingAds, _ := storage.LoadAds()
	if existingAds == nil {
		existingAds = &models.AdDeliveryResponse{
			Ads:       []models.Ad{},
			UpdatedAt: time.Now(),
		}
	}

	// Add new ad
	newAd := models.Ad{
		ID:         adID,
		Title:      title,
		Type:       "image",
		ContentURL: "file://" + destFile,
		Duration:   duration,
		Priority:   1,
	}

	// Check if ad already exists, replace it
	found := false
	for i, ad := range existingAds.Ads {
		if ad.ID == adID {
			existingAds.Ads[i] = newAd
			found = true
			break
		}
	}

	if !found {
		existingAds.Ads = append(existingAds.Ads, newAd)
	}

	existingAds.UpdatedAt = time.Now()

	if err := storage.SaveAds(existingAds); err != nil {
		log.Fatalf("Failed to save ads: %v", err)
	}

	fmt.Printf("[OK] Added image ad: %s\n", adID)
	fmt.Printf("   Title: %s\n", title)
	fmt.Printf("   Duration: %d seconds\n", duration)
	fmt.Printf("   File: %s\n", destFile)
}

func addVideoAd(storage *ads.Storage, adID, videoPath, title string, duration int) {
	// Check if video file exists
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		log.Fatalf("Video file not found: %s", videoPath)
	}

	// Copy video to media directory
	mediaDir := storage.GetMediaDir()
	adMediaDir := filepath.Join(mediaDir, adID)
	os.MkdirAll(adMediaDir, 0755)

	ext := filepath.Ext(videoPath)
	destFile := filepath.Join(adMediaDir, adID+ext)

	// Read and copy file
	data, err := os.ReadFile(videoPath)
	if err != nil {
		log.Fatalf("Failed to read video file: %v", err)
	}

	if err := os.WriteFile(destFile, data, 0644); err != nil {
		log.Fatalf("Failed to copy video file: %v", err)
	}

	// Load existing ads
	existingAds, _ := storage.LoadAds()
	if existingAds == nil {
		existingAds = &models.AdDeliveryResponse{
			Ads:       []models.Ad{},
			UpdatedAt: time.Now(),
		}
	}

	// Add new ad
	newAd := models.Ad{
		ID:         adID,
		Title:      title,
		Type:       "video",
		ContentURL: "file://" + destFile,
		Duration:   duration,
		Priority:   1,
	}

	// Check if ad already exists, replace it
	found := false
	for i, ad := range existingAds.Ads {
		if ad.ID == adID {
			existingAds.Ads[i] = newAd
			found = true
			break
		}
	}

	if !found {
		existingAds.Ads = append(existingAds.Ads, newAd)
	}

	existingAds.UpdatedAt = time.Now()

	if err := storage.SaveAds(existingAds); err != nil {
		log.Fatalf("Failed to save ads: %v", err)
	}

	fmt.Printf("[OK] Added video ad: %s\n", adID)
	fmt.Printf("   Title: %s\n", title)
	fmt.Printf("   Duration: %d seconds\n", duration)
	fmt.Printf("   File: %s\n", destFile)
}

func addTextAd(storage *ads.Storage, adID, textContent, title string, duration int) {
	// Save text to file
	mediaDir := storage.GetMediaDir()
	adMediaDir := filepath.Join(mediaDir, adID)
	os.MkdirAll(adMediaDir, 0755)

	textFile := filepath.Join(adMediaDir, adID+".txt")
	if err := os.WriteFile(textFile, []byte(textContent), 0644); err != nil {
		log.Fatalf("Failed to create text file: %v", err)
	}

	// Load existing ads
	existingAds, _ := storage.LoadAds()
	if existingAds == nil {
		existingAds = &models.AdDeliveryResponse{
			Ads:       []models.Ad{},
			UpdatedAt: time.Now(),
		}
	}

	// Add new ad
	newAd := models.Ad{
		ID:         adID,
		Title:      title,
		Type:       "text",
		ContentURL: "file://" + textFile,
		Duration:   duration,
		Priority:   1,
	}

	// Check if ad already exists, replace it
	found := false
	for i, ad := range existingAds.Ads {
		if ad.ID == adID {
			existingAds.Ads[i] = newAd
			found = true
			break
		}
	}

	if !found {
		existingAds.Ads = append(existingAds.Ads, newAd)
	}

	existingAds.UpdatedAt = time.Now()

	if err := storage.SaveAds(existingAds); err != nil {
		log.Fatalf("Failed to save ads: %v", err)
	}

	fmt.Printf("[OK] Added text ad: %s\n", adID)
	fmt.Printf("   Title: %s\n", title)
	fmt.Printf("   Duration: %d seconds\n", duration)
	fmt.Printf("   File: %s\n", textFile)
}

func listAds(storage *ads.Storage) {
	ads, err := storage.LoadAds()
	if err != nil {
		fmt.Println("[INFO] No ads found")
		return
	}

	fmt.Printf("Found %d ads:\n\n", len(ads.Ads))
	for i, ad := range ads.Ads {
		fmt.Printf("%d. ID: %s\n", i+1, ad.ID)
		fmt.Printf("   Title: %s\n", ad.Title)
		fmt.Printf("   Type: %s\n", ad.Type)
		fmt.Printf("   Duration: %d seconds\n", ad.Duration)
		fmt.Printf("   Priority: %d\n", ad.Priority)
		fmt.Printf("   Content URL: %s\n", ad.ContentURL)
		fmt.Println()
	}
}

func clearAds(storage *ads.Storage) {
	emptyAds := &models.AdDeliveryResponse{
		Ads:       []models.Ad{},
		UpdatedAt: time.Now(),
	}

	if err := storage.SaveAds(emptyAds); err != nil {
		log.Fatalf("Failed to clear ads: %v", err)
	}

	fmt.Println("[OK] Cleared all ads")
}

func getOptionalArg(args []string, index int, defaultValue string) string {
	if index < len(args) {
		return args[index]
	}
	return defaultValue
}

func getOptionalIntArg(args []string, index int, defaultValue int) int {
	if index < len(args) {
		var val int
		if _, err := fmt.Sscanf(args[index], "%d", &val); err == nil {
			return val
		}
	}
	return defaultValue
}

