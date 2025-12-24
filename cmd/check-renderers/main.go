package main

import (
	"fmt"
	"mnemoCast-client/internal/player/renderers"
	"os/exec"
)

func main() {
	fmt.Println("Checking available renderers...")
	fmt.Println()

	// Check image viewers
	fmt.Println("Image Viewers:")
	checkCommand("feh")
	checkCommand("imv")
	checkCommand("sxiv")
	checkCommand("xdg-open")
	fmt.Println()

	// Check video players
	fmt.Println("Video Players:")
	checkCommand("mpv")
	checkCommand("vlc")
	checkCommand("ffplay")
	fmt.Println()

	// Check browsers
	fmt.Println("Browsers:")
	checkCommand("firefox")
	checkCommand("chromium")
	checkCommand("chrome")
	fmt.Println()

	// Test renderers
	fmt.Println("Renderer Status:")
	imgRenderer := renderers.NewImageRenderer()
	fmt.Printf("  Image Renderer: %s\n", func() string {
		if imgRenderer != nil {
			return "Available"
		}
		return "Not Available"
	}())

	vidRenderer := renderers.NewVideoRenderer()
	fmt.Printf("  Video Renderer: %s\n", func() string {
		if vidRenderer != nil {
			return "Available"
		}
		return "Not Available"
	}())

	htmlRenderer := renderers.NewHTMLRenderer()
	fmt.Printf("  HTML Renderer: %s\n", func() string {
		if htmlRenderer != nil {
			return "Available"
		}
		return "Not Available"
	}())

	textRenderer := renderers.NewTextRenderer()
	fmt.Printf("  Text Renderer: %s\n", func() string {
		if textRenderer != nil {
			return "Available"
		}
		return "Not Available"
	}())
}

func checkCommand(cmd string) {
	path, err := exec.LookPath(cmd)
	if err == nil {
		fmt.Printf("  ✓ %s: %s\n", cmd, path)
	} else {
		fmt.Printf("  ✗ %s: Not found\n", cmd)
	}
}

