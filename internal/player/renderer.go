package player

import (
	"fmt"
	"mnemoCast-client/internal/models"
	"mnemoCast-client/internal/player/renderers"
)

// RendererStatus represents the status of a renderer (aliased from renderers package)
type RendererStatus = renderers.RendererStatus

// Renderer defines the interface for ad renderers
type Renderer interface {
	// CanRender checks if this renderer can handle the given ad type
	CanRender(ad *models.Ad) bool
	
	// Render displays the ad using the local file path
	Render(ad *models.Ad, localPath string) error
	
	// Stop stops the current rendering
	Stop() error
	
	// GetStatus returns the current renderer status
	GetStatus() RendererStatus
}

// RendererManager manages multiple renderers and selects the appropriate one
type RendererManager struct {
	renderers []Renderer
	current   Renderer
}

// NewRendererManager creates a new renderer manager
func NewRendererManager() *RendererManager {
	return &RendererManager{
		renderers: []Renderer{
			renderers.NewImageRenderer(),
			renderers.NewVideoRenderer(),
			renderers.NewHTMLRenderer(),
			renderers.NewTextRenderer(),
		},
	}
}

// GetRenderer returns the appropriate renderer for an ad
func (rm *RendererManager) GetRenderer(ad *models.Ad) Renderer {
	for _, renderer := range rm.renderers {
		if renderer.CanRender(ad) {
			return renderer
		}
	}
	return nil // No renderer found
}

// Render renders an ad using the appropriate renderer
func (rm *RendererManager) Render(ad *models.Ad, localPath string) error {
	// Stop current renderer if any
	if rm.current != nil {
		rm.current.Stop()
	}
	
	// Get appropriate renderer
	renderer := rm.GetRenderer(ad)
	if renderer == nil {
		return &RendererError{
			AdID:  ad.ID,
			AdType: ad.Type,
			Message: "no renderer available for ad type",
		}
	}
	
	rm.current = renderer
	return renderer.Render(ad, localPath)
}

// Stop stops the current renderer
func (rm *RendererManager) Stop() error {
	if rm.current != nil {
		err := rm.current.Stop()
		rm.current = nil
		return err
	}
	return nil
}

// RendererError represents an error in rendering
type RendererError struct {
	AdID    string
	AdType  string
	Message string
}

func (e *RendererError) Error() string {
	return fmt.Sprintf("renderer error for ad %s (type: %s): %s", e.AdID, e.AdType, e.Message)
}

