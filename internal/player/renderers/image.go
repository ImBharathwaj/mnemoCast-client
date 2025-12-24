package renderers

import (
	"fmt"
	"log"
	"mnemoCast-client/internal/models"
	"os"
	"os/exec"
	"time"
)

// ImageRenderer renders image ads
type ImageRenderer struct {
	currentProcess *os.Process
	displayCmd     string
	status         RendererStatus
}

// NewImageRenderer creates a new image renderer
func NewImageRenderer() *ImageRenderer {
	// Try to find an image viewer command
	displayCmd := findImageCommand()
	
	return &ImageRenderer{
		displayCmd: displayCmd,
		status: RendererStatus{
			IsPlaying: false,
		},
	}
}

// CanRender checks if this renderer can handle the ad type
func (r *ImageRenderer) CanRender(ad *models.Ad) bool {
	adType := ad.Type
	return adType == "image" || adType == "jpg" || adType == "jpeg" || 
		   adType == "png" || adType == "gif" || adType == "webp"
}

// Render displays the image ad
func (r *ImageRenderer) Render(ad *models.Ad, localPath string) error {
	// Stop any existing process
	r.Stop()
	
	if r.displayCmd == "" {
		return fmt.Errorf("no image viewer available")
	}
	
	// Check if file exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %s", localPath)
	}
	
	log.Printf("[%s] [RENDER] Rendering image ad: %s using %s", 
		time.Now().Format("15:04:05.000"), ad.ID, r.displayCmd)
	
	// Launch image viewer in fullscreen mode
	var cmd *exec.Cmd
	switch r.displayCmd {
	case "feh":
		cmd = exec.Command("feh", "--fullscreen", "--auto-zoom", localPath)
	case "imv":
		cmd = exec.Command("imv", "-f", localPath)
	case "sxiv":
		cmd = exec.Command("sxiv", "-f", localPath)
	default:
		// Fallback: try to open with system default
		cmd = exec.Command("xdg-open", localPath)
	}
	
	// Start process in background
	log.Printf("[%s] [RENDER] Executing command: %s %v", 
		time.Now().Format("15:04:05.000"), cmd.Path, cmd.Args)
	
	if err := cmd.Start(); err != nil {
		log.Printf("[%s] [RENDER] [ERROR] Failed to start image viewer: %v", 
			time.Now().Format("15:04:05.000"), err)
		return fmt.Errorf("failed to start image viewer: %w", err)
	}
	
	r.currentProcess = cmd.Process
	r.status.IsPlaying = true
	r.status.Error = nil
	
	log.Printf("[%s] [RENDER] [OK] Image viewer started (PID: %d)", 
		time.Now().Format("15:04:05.000"), cmd.Process.Pid)
	
	return nil
}

// Stop stops the image rendering
func (r *ImageRenderer) Stop() error {
	if r.currentProcess != nil {
		// Try to kill the process
		if err := r.currentProcess.Kill(); err != nil {
			log.Printf("[%s] [RENDER] Failed to stop image viewer: %v", 
				time.Now().Format("15:04:05.000"), err)
		}
		r.currentProcess = nil
	}
	
	r.status.IsPlaying = false
	return nil
}

// GetStatus returns the current renderer status
func (r *ImageRenderer) GetStatus() RendererStatus {
	return r.status
}

// findImageCommand finds an available image viewer command
func findImageCommand() string {
	commands := []string{"feh", "imv", "sxiv", "xdg-open"}
	
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err == nil {
			return cmd
		}
	}
	
	return "" // No image viewer found
}

