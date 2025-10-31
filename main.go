package main

import (
	"flag"
	"log"
	"os"

	"github.com/clobrano/mcp-desktop-notification/internal/config"
	"github.com/clobrano/mcp-desktop-notification/internal/mcp"
	"github.com/clobrano/mcp-desktop-notification/internal/notifier"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "", "Path to configuration file (default: platform-specific)")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	dryRun := flag.Bool("dry-run", false, "Dry run mode (log notifications without sending)")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.LoadConfig(*configPath)
	} else {
		cfg, err = config.LoadDefaultConfig()
	}

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with command-line flags
	if *verbose {
		cfg.Notification.Verbose = true
	}
	if *dryRun {
		cfg.Notification.DryRun = true
	}

	// Log configuration if verbose
	if cfg.Notification.Verbose {
		log.Printf("[Main] Configuration loaded - DryRun: %v, Verbose: %v",
			cfg.Notification.DryRun, cfg.Notification.Verbose)
	}

	// Create notifier
	noti, err := notifier.NewNotifier(cfg)
	if err != nil {
		log.Fatalf("Failed to create notifier: %v", err)
	}

	if cfg.Notification.Verbose {
		log.Printf("[Main] Notifier created successfully")
	}

	// Create and start MCP server
	server := mcp.NewServer(cfg, noti)

	if cfg.Notification.Verbose {
		log.Printf("[Main] Starting MCP server...")
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	os.Exit(0)
}
