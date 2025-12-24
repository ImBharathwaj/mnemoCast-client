package renderers

import (
	"fmt"
	"log"
	"mnemoCast-client/internal/models"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// HTMLRenderer renders HTML ads using a local web server
type HTMLRenderer struct {
	server     *http.Server
	serverPort int
	browserCmd string
	status     RendererStatus
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	// Try to find a browser command
	browserCmd := findBrowserCommand()
	
	return &HTMLRenderer{
		browserCmd: browserCmd,
		status: RendererStatus{
			IsPlaying: false,
		},
	}
}

// CanRender checks if this renderer can handle the ad type
func (r *HTMLRenderer) CanRender(ad *models.Ad) bool {
	return ad.Type == "html"
}

// Render displays the HTML ad
func (r *HTMLRenderer) Render(ad *models.Ad, localPath string) error {
	// Stop any existing rendering
	r.Stop()
	
	if r.browserCmd == "" {
		return fmt.Errorf("no browser available")
	}
	
	// Check if file exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("HTML file not found: %s", localPath)
	}
	
	log.Printf("[%s] [RENDER] Rendering HTML ad: %s", time.Now().Format("15:04:05.000"), ad.ID)
	
	// Start local HTTP server
	port, err := r.startServer(localPath)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	r.serverPort = port
	
	// Open browser
	url := fmt.Sprintf("http://localhost:%d", port)
	if err := r.openBrowser(url); err != nil {
		r.Stop()
		return fmt.Errorf("failed to open browser: %w", err)
	}
	
	r.status.IsPlaying = true
	r.status.Error = nil
	
	return nil
}

// Stop stops the HTML rendering
func (r *HTMLRenderer) Stop() error {
	if r.server != nil {
		if err := r.server.Close(); err != nil {
			log.Printf("[%s] [RENDER] Failed to stop HTTP server: %v", 
				time.Now().Format("15:04:05.000"), err)
		}
		r.server = nil
	}
	
	r.status.IsPlaying = false
	return nil
}

// GetStatus returns the current renderer status
func (r *HTMLRenderer) GetStatus() RendererStatus {
	return r.status
}

// startServer starts a local HTTP server to serve the HTML file
func (r *HTMLRenderer) startServer(filePath string) (int, error) {
	// Find an available port
	port := 8081 // Start from 8081
	
	dir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath.Join(dir, fileName))
	})
	
	r.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	
	// Start server in goroutine
	go func() {
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[%s] [RENDER] HTTP server error: %v", time.Now().Format("15:04:05.000"), err)
		}
	}()
	
	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	return port, nil
}

// openBrowser opens the URL in the default browser
func (r *HTMLRenderer) openBrowser(url string) error {
	var cmd *exec.Cmd
	
	switch r.browserCmd {
	case "firefox":
		cmd = exec.Command("firefox", "--kiosk", url)
	case "chromium", "chrome":
		cmd = exec.Command(r.browserCmd, "--kiosk", url)
	default:
		// Fallback: use xdg-open
		cmd = exec.Command("xdg-open", url)
	}
	
	return cmd.Start()
}

// findBrowserCommand finds an available browser command
func findBrowserCommand() string {
	commands := []string{"firefox", "chromium", "chrome", "xdg-open"}
	
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err == nil {
			return cmd
		}
	}
	
	return "" // No browser found
}

