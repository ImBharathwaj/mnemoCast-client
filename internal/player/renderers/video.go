package renderers

import (
	"fmt"
	"log"
	"mnemoCast-client/internal/models"
	"os"
	"os/exec"
	"time"
)

// VideoRenderer renders video ads
type VideoRenderer struct {
	currentProcess *os.Process
	playerCmd      string
	status         RendererStatus
}

// NewVideoRenderer creates a new video renderer
func NewVideoRenderer() *VideoRenderer {
	// Try to find a video player command
	playerCmd := findVideoCommand()
	
	return &VideoRenderer{
		playerCmd: playerCmd,
		status: RendererStatus{
			IsPlaying: false,
		},
	}
}

// CanRender checks if this renderer can handle the ad type
func (r *VideoRenderer) CanRender(ad *models.Ad) bool {
	adType := ad.Type
	return adType == "video" || adType == "mp4" || adType == "webm" || 
		   adType == "mov" || adType == "avi"
}

// Render displays the video ad
func (r *VideoRenderer) Render(ad *models.Ad, localPath string) error {
	// Stop any existing process
	r.Stop()
	
	if r.playerCmd == "" {
		return fmt.Errorf("no video player available")
	}
	
	// Check if file exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("video file not found: %s", localPath)
	}
	
	log.Printf("[%s] [RENDER] Rendering video ad: %s using %s", 
		time.Now().Format("15:04:05.000"), ad.ID, r.playerCmd)
	
	// Launch video player in fullscreen mode
	var cmd *exec.Cmd
	switch r.playerCmd {
	case "mpv":
		cmd = exec.Command("mpv", "--fullscreen", "--loop=no", localPath)
	case "vlc":
		cmd = exec.Command("vlc", "--fullscreen", "--no-loop", localPath)
	case "ffplay":
		cmd = exec.Command("ffplay", "-fs", "-autoexit", localPath)
	default:
		// Fallback: try to open with system default
		cmd = exec.Command("xdg-open", localPath)
	}
	
	// Start process in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start video player: %w", err)
	}
	
	r.currentProcess = cmd.Process
	r.status.IsPlaying = true
	r.status.Error = nil
	
	return nil
}

// Stop stops the video rendering
func (r *VideoRenderer) Stop() error {
	if r.currentProcess != nil {
		// Try to kill the process
		if err := r.currentProcess.Kill(); err != nil {
			log.Printf("[%s] [RENDER] Failed to stop video player: %v", 
				time.Now().Format("15:04:05.000"), err)
		}
		r.currentProcess = nil
	}
	
	r.status.IsPlaying = false
	return nil
}

// GetStatus returns the current renderer status
func (r *VideoRenderer) GetStatus() RendererStatus {
	return r.status
}

// findVideoCommand finds an available video player command
func findVideoCommand() string {
	commands := []string{"mpv", "vlc", "ffplay", "xdg-open"}
	
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err == nil {
			return cmd
		}
	}
	
	return "" // No video player found
}

