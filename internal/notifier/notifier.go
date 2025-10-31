package notifier

import (
	"fmt"
	"log"
	"runtime"

	"github.com/clobrano/mcp-desktop-notification/internal/config"
	"github.com/gen2brain/beeep"
)

// Notifier is the interface for sending notifications
type Notifier interface {
	Send(title, message, level string) error
}

// LibraryNotifier sends notifications using the beeep library
type LibraryNotifier struct {
	config *config.Config
}

// DryRunNotifier logs notifications without sending them
type DryRunNotifier struct {
	config *config.Config
}

// NewNotifier creates a new notifier based on configuration
func NewNotifier(cfg *config.Config) (Notifier, error) {
	// If dry run mode is enabled, return dry run notifier
	if cfg.Notification.DryRun {
		return &DryRunNotifier{config: cfg}, nil
	}

	// Create library-based notifier
	if cfg.Notification.Mode == "library" {
		return &LibraryNotifier{config: cfg}, nil
	}

	return nil, fmt.Errorf("unsupported notification mode: %s", cfg.Notification.Mode)
}

// Send sends a notification using the beeep library
func (n *LibraryNotifier) Send(title, message, level string) error {
	// Get icon based on level
	icon := n.getIcon(level)

	// Log if verbose
	if n.config.Notification.Verbose {
		log.Printf("[LibraryNotifier] Sending notification - Title: %s, Message: %s, Level: %s, Icon: %s, Platform: %s",
			title, message, level, icon, runtime.GOOS)
	}

	// Send notification using beeep
	err := beeep.Notify(title, message, icon)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

// Send logs the notification without sending it
func (n *DryRunNotifier) Send(title, message, level string) error {
	log.Printf("[DRY RUN] Would send notification - Title: %s, Message: %s, Level: %s, Platform: %s",
		title, message, level, runtime.GOOS)
	return nil
}

// getUrgency returns the urgency level for a notification level
func (n *LibraryNotifier) getUrgency(level string) string {
	if levelConfig, ok := n.config.Notification.Levels[level]; ok {
		return levelConfig.Urgency
	}
	return "normal" // default
}

// getIcon returns the icon for a notification level
func (n *LibraryNotifier) getIcon(level string) string {
	if levelConfig, ok := n.config.Notification.Levels[level]; ok {
		return levelConfig.Icon
	}
	return "" // default (no icon)
}
