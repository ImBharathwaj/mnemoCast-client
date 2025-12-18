package main

import (
	"bufio"
	"fmt"
	"log"
	"mnemoCast-client/internal/ads"
	"mnemoCast-client/internal/client"
	"mnemoCast-client/internal/config"
	"mnemoCast-client/internal/credentials"
	"mnemoCast-client/internal/heartbeat"
	"mnemoCast-client/internal/identity"
	"mnemoCast-client/internal/models"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// Set up config directory
	configDir := filepath.Join(homeDir, ".mnemocast")
	
	fmt.Println("MnemoCast Screen System")
	fmt.Println("============================")
	fmt.Printf("Config directory: %s\n\n", configDir)

	// Initialize identity manager
	identityManager := identity.NewManager(configDir)
	
	// Try to load existing identity (optional - will be loaded from server after connection)
	var screenIdentity *models.ScreenIdentity
	fmt.Println("Loading screen identity...")
	loadedIdentity, err := identityManager.LoadIdentity()
	if err == nil {
		screenIdentity = loadedIdentity
		fmt.Printf("[OK] Screen ID: %s\n", screenIdentity.ID)
		fmt.Printf("   Name: %s\n", screenIdentity.Name)
		
		// Display location
		locationParts := []string{}
		if screenIdentity.Country != "" && screenIdentity.Country != "Unknown" {
			locationParts = append(locationParts, screenIdentity.Country)
		}
		if screenIdentity.City != "" && screenIdentity.City != "Unknown" {
			locationParts = append(locationParts, screenIdentity.City)
		}
		if screenIdentity.Area != "" && screenIdentity.Area != "Unknown" {
			locationParts = append(locationParts, screenIdentity.Area)
		}
		if len(locationParts) > 0 {
			fmt.Printf("   Location: %s", strings.Join(locationParts, ", "))
			if screenIdentity.VenueType != "" && screenIdentity.VenueType != "unknown" {
				fmt.Printf(" (%s)", screenIdentity.VenueType)
			}
			fmt.Println()
		}
		fmt.Println()
	} else {
		fmt.Println("[WARN] Screen identity not found (will be loaded from server after connection)")
		fmt.Println()
	}

	// Initialize config loader
	configLoader := config.NewLoader(configDir)
	
	// Load configuration
	fmt.Println("Loading configuration...")
	screenConfig, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("[OK] Ad Server URL: %s\n", screenConfig.AdServerURL)
	fmt.Printf("   Heartbeat Interval: %d seconds\n", screenConfig.HeartbeatInterval)
	if screenConfig.AdFetchInterval > 0 {
		fmt.Printf("   Ad Fetch Interval: %d seconds\n", screenConfig.AdFetchInterval)
	}
	fmt.Printf("   Retry Attempts: %d\n", screenConfig.RetryAttempts)
	fmt.Printf("   Retry Delay: %d seconds\n", screenConfig.RetryDelay)
	fmt.Println()

	// Update identity in config if needed
	if screenConfig.Identity.ID == "" {
		screenConfig.Identity = *screenIdentity
		if err := configLoader.Save(screenConfig); err != nil {
			log.Printf("Warning: Failed to save updated config: %v", err)
		}
	}

	// Initialize credentials manager
	credManager := credentials.NewManager(configDir)
	
	// Check credentials
	fmt.Println("Checking credentials...")
	var screenID, passkey string
	if credManager.Exists() {
		loadedScreenID, loadedPasskey, err := credManager.GetCredentials()
		if err != nil {
			log.Printf("Warning: Failed to load credentials: %v", err)
		} else {
			screenID = loadedScreenID
			passkey = loadedPasskey
			fmt.Println("[OK] Credentials: Configured")
			// Mask passkey for display
			if len(passkey) > 8 {
				masked := passkey[:4] + "..." + passkey[len(passkey)-4:]
				fmt.Printf("   Screen ID: %s\n", screenID)
				fmt.Printf("   Passkey: %s\n", masked)
			} else {
				fmt.Printf("   Screen ID: %s\n", screenID)
				fmt.Println("   Passkey: [hidden]")
			}
		}
	} else {
		fmt.Println("[WARN] Credentials: Not configured")
		fmt.Println()
		fmt.Println("INFO: To configure credentials:")
		fmt.Println("   1. Register a screen on the server")
		fmt.Println("   2. Get the Screen ID and Passkey from the server")
		fmt.Println("   3. Enter them below")
		fmt.Println()
		fmt.Print("Would you like to configure credentials now? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		
		if response == "y" || response == "yes" {
			fmt.Print("Enter Screen ID (from server): ")
			inputScreenID, _ := reader.ReadString('\n')
			inputScreenID = strings.TrimSpace(inputScreenID)
			
			fmt.Print("Enter Passkey (from server): ")
			inputPasskey, _ := reader.ReadString('\n')
			inputPasskey = strings.TrimSpace(inputPasskey)
			
			if inputScreenID != "" && inputPasskey != "" {
				if err := credManager.SetCredentials(inputScreenID, inputPasskey); err != nil {
					log.Printf("Failed to save credentials: %v", err)
				} else {
					screenID = inputScreenID
					passkey = inputPasskey
					fmt.Println("[OK] Credentials saved successfully!")
				}
			} else {
				fmt.Println("[WARN] Both Screen ID and Passkey are required")
			}
		}
	}
	fmt.Println()

	// Initialize ad server client, heartbeat, and ad fetcher if credentials exist
	var heartbeatScheduler *heartbeat.Scheduler
	var adFetcher *ads.Fetcher
	var adClient *client.Client
	
	if screenID != "" && passkey != "" {
		// Create ad server client with screen ID and passkey
		adClient = client.NewClient(screenConfig.AdServerURL, screenID, passkey)
		
		// Try to connect/authenticate with server (optional - heartbeat will work even if this fails)
		fmt.Println("Connecting to ad server...")
		fmt.Println("   Authenticating with server...")
		connectedScreen, err := adClient.Connect()
		if err != nil {
			log.Printf("[WARN] Connection failed: %v", err)
			fmt.Println("   [WARN] Could not connect to ad server")
			fmt.Println("   Note: Heartbeat will still be attempted...")
		} else {
			fmt.Printf("   [OK] Connected successfully!\n")
			fmt.Printf("      Screen ID: %s\n", connectedScreen.ID)
			fmt.Printf("      Name: %s\n", connectedScreen.Name)
			if connectedScreen.City != "" || connectedScreen.Area != "" {
				location := []string{}
				if connectedScreen.City != "" {
					location = append(location, connectedScreen.City)
				}
				if connectedScreen.Area != "" {
					location = append(location, connectedScreen.Area)
				}
				fmt.Printf("      Location: %s\n", strings.Join(location, ", "))
			}
			fmt.Printf("      Status: %s\n", func() string {
				if connectedScreen.IsOnline {
					return "Online"
				}
				return "Offline"
			}())
			
			// Update identity from server response
			if err := identityManager.UpdateIdentityFromServer(connectedScreen); err != nil {
				log.Printf("Warning: Failed to update identity: %v", err)
			} else {
				screenIdentity, _ = identityManager.LoadIdentity()
			}
		}
		
		// Test heartbeat
		fmt.Println()
		fmt.Println("   Testing heartbeat...")
		if err := adClient.Heartbeat(screenID); err != nil {
			log.Printf("[WARN] Heartbeat test failed: %v", err)
			fmt.Println("   [WARN] Heartbeat test failed")
		} else {
			fmt.Println("   [OK] Heartbeat successful!")
		}

		// Start heartbeat scheduler (even if connection failed, heartbeat might still work)
		fmt.Println()
		fmt.Println("Starting heartbeat scheduler...")
		heartbeatScheduler = heartbeat.NewScheduler(
			adClient,
			identityManager,
			screenID,
			screenConfig.HeartbeatInterval,
			screenConfig.RetryAttempts,
			screenConfig.RetryDelay,
		)
		heartbeatScheduler.Start()
		fmt.Printf("   [OK] Heartbeat scheduler started (interval: %d seconds)\n", screenConfig.HeartbeatInterval)

		// Start ad fetcher
		if screenConfig.AdFetchInterval > 0 {
			fmt.Println()
			fmt.Println("Starting ad fetcher...")
			adFetcher = ads.NewFetcher(
				adClient,
				screenID,
				configDir,
				screenConfig.AdFetchInterval,
				screenConfig.RetryAttempts,
				screenConfig.RetryDelay,
			)
			
			// Try to load existing ads from storage
			if storedAds, err := adFetcher.LoadAdsFromStorage(); err == nil {
				fmt.Printf("   [INFO] Loaded %d ads from storage\n", len(storedAds.Ads))
				if len(storedAds.Ads) > 0 {
					fmt.Printf("   [INFO] Ads storage location: %s\n", adFetcher.GetStorage().GetAdsDir())
				}
			}
			
			adFetcher.Start()
			fmt.Printf("   [OK] Ad fetcher started (interval: %d seconds)\n", screenConfig.AdFetchInterval)
			fmt.Printf("   [INFO] Ads will be stored in: %s/ads/\n", configDir)
		}
	}

	fmt.Println()
	fmt.Println("[OK] Screen system initialized successfully!")
	fmt.Println()
	
	// If heartbeat is running, show status and wait for interrupt
	if heartbeatScheduler != nil {
		fmt.Println("Heartbeat system is running...")
		fmt.Println("   Press Ctrl+C to stop")
		fmt.Println()
		
		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		
		// Status update ticker
		statusTicker := time.NewTicker(30 * time.Second)
		defer statusTicker.Stop()
		
		// Main loop
		for {
			select {
			case <-sigChan:
				fmt.Println()
				fmt.Println("Shutting down...")
				heartbeatScheduler.Stop()
				if adFetcher != nil {
					adFetcher.Stop()
				}
				fmt.Println("[OK] Shutdown complete")
				return
			case <-statusTicker.C:
				stats := heartbeatScheduler.GetStats()
				status := stats["status"].(string)
				connected := stats["connected"].(bool)
				
				statusIcon := "[ERROR]"
				if connected {
					statusIcon = "[OK]"
				}
				
				fmt.Printf("%s Heartbeat Status: %s", statusIcon, status)
				if lastSent, ok := stats["timeSinceLastSent"]; ok {
					fmt.Printf(" (last sent: %s ago)", lastSent)
				}
				fmt.Println()
			}
		}
	} else {
		// System is running but waiting for credentials
		fmt.Println("[WARN] System is running but heartbeat is not active")
		fmt.Println("   Credentials are required to connect to ad server")
		fmt.Println()
		fmt.Println("INFO: To enable heartbeat system:")
		fmt.Println("   1. Register a screen on the server")
		fmt.Println("   2. Get the Screen ID and Passkey from the server")
		fmt.Println("   3. Stop this process (Ctrl+C)")
		fmt.Println("   4. Run again: ./bin/screen")
		fmt.Println("   5. Enter 'y' when prompted and configure credentials")
		fmt.Println()
		fmt.Println("Waiting for credentials...")
		fmt.Println("   Press Ctrl+C to exit")
		fmt.Println()
		
		// Wait for interrupt
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		
		fmt.Println()
		fmt.Println("Shutting down...")
		fmt.Println("[OK] Shutdown complete")
		return
	}
}

