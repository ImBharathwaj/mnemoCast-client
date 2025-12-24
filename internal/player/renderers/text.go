package renderers

import (
	"fmt"
	"log"
	"mnemoCast-client/internal/models"
	"os"
	"strings"
	"time"
)

// TextRenderer renders text ads in the terminal
type TextRenderer struct {
	status RendererStatus
}

// NewTextRenderer creates a new text renderer
func NewTextRenderer() *TextRenderer {
	return &TextRenderer{
		status: RendererStatus{
			IsPlaying: false,
		},
	}
}

// CanRender checks if this renderer can handle the ad type
func (r *TextRenderer) CanRender(ad *models.Ad) bool {
	return ad.Type == "text"
}

// Render displays the text ad
func (r *TextRenderer) Render(ad *models.Ad, localPath string) error {
	// Stop any existing rendering
	r.Stop()
	
	log.Printf("[%s] [RENDER] Rendering text ad: %s", time.Now().Format("15:04:05.000"), ad.ID)
	
	// Read text content
	var content string
	if localPath != "" {
		// Try to read from file
		if data, err := os.ReadFile(localPath); err == nil {
			content = string(data)
		} else {
			// Fallback to ad title or ID
			content = ad.Title
			if content == "" {
				content = ad.ID
			}
		}
	} else {
		// Use ad title or ID
		content = ad.Title
		if content == "" {
			content = ad.ID
		}
	}
	
	// Display text in terminal with formatting
	log.Printf("[%s] [RENDER] Displaying text ad: %s", time.Now().Format("15:04:05.000"), ad.ID)
	r.displayText(ad, content)
	
	r.status.IsPlaying = true
	r.status.Error = nil
	
	log.Printf("[%s] [RENDER] [OK] Text ad displayed", time.Now().Format("15:04:05.000"))
	
	return nil
}

// Stop stops the text rendering
func (r *TextRenderer) Stop() error {
	// Clear terminal (optional)
	r.status.IsPlaying = false
	return nil
}

// GetStatus returns the current renderer status
func (r *TextRenderer) GetStatus() RendererStatus {
	return r.status
}

// displayText displays text in the terminal with formatting
func (r *TextRenderer) displayText(ad *models.Ad, content string) {
	// Simple terminal display
	fmt.Println("\n" + "=" + strings.Repeat("=", 78) + "=")
	fmt.Printf("TEXT AD: %s\n", ad.ID)
	if ad.Title != "" {
		fmt.Printf("Title: %s\n", ad.Title)
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(content)
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

