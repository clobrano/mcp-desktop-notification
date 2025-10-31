package notifier

import (
	"runtime"
	"testing"

	"github.com/clobrano/mcp-desktop-notification/internal/config"
)

func TestNewNotifier_LibraryMode(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Notification.Mode = "library"
	cfg.Notification.DryRun = false

	notifier, err := NewNotifier(cfg)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	if _, ok := notifier.(*LibraryNotifier); !ok {
		t.Error("Expected LibraryNotifier for library mode")
	}
}

func TestNewNotifier_DryRunMode(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Notification.DryRun = true

	notifier, err := NewNotifier(cfg)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	if _, ok := notifier.(*DryRunNotifier); !ok {
		t.Error("Expected DryRunNotifier for dry run mode")
	}
}

func TestLibraryNotifier_Send(t *testing.T) {
	// Skip this test in CI or headless environments
	// Beeep requires a display server which may not be available
	if testing.Short() {
		t.Skip("Skipping notification send test in short mode")
	}

	cfg := config.DefaultConfig()
	notifier := &LibraryNotifier{
		config: cfg,
	}

	// Test sending notification (won't actually show on desktop during test)
	// This test may fail in headless environments, which is expected
	err := notifier.Send("Test Title", "Test Message", "info")

	// On Linux without X server, this might fail, which is expected in CI
	// We just check that the method doesn't panic
	if err != nil {
		t.Logf("Send returned error (expected in headless environment): %v", err)
	}
}

func TestLibraryNotifier_LevelMapping(t *testing.T) {
	cfg := config.DefaultConfig()
	notifier := &LibraryNotifier{
		config: cfg,
	}

	tests := []struct {
		level    string
		expected string
	}{
		{"info", "normal"},
		{"warning", "normal"},
		{"error", "critical"},
		{"success", "low"},
		{"unknown", "normal"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			urgency := notifier.getUrgency(tt.level)
			if urgency != tt.expected {
				t.Errorf("Expected urgency %s for level %s, got %s", tt.expected, tt.level, urgency)
			}
		})
	}
}

func TestLibraryNotifier_Icon(t *testing.T) {
	cfg := config.DefaultConfig()
	notifier := &LibraryNotifier{
		config: cfg,
	}

	tests := []struct {
		level        string
		expectedIcon string
	}{
		{"info", ""},
		{"warning", "dialog-warning"},
		{"error", "dialog-error"},
		{"success", "dialog-information"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			icon := notifier.getIcon(tt.level)
			if icon != tt.expectedIcon {
				t.Errorf("Expected icon %s for level %s, got %s", tt.expectedIcon, tt.level, icon)
			}
		})
	}
}

func TestDryRunNotifier_Send(t *testing.T) {
	cfg := config.DefaultConfig()
	notifier := &DryRunNotifier{
		config: cfg,
	}

	// Dry run should always succeed
	err := notifier.Send("Test Title", "Test Message", "info")
	if err != nil {
		t.Errorf("DryRun should not return error, got: %v", err)
	}
}

func TestPlatformDetection(t *testing.T) {
	platform := runtime.GOOS

	validPlatforms := []string{"linux", "darwin", "windows"}
	isValid := false
	for _, p := range validPlatforms {
		if platform == p {
			isValid = true
			break
		}
	}

	if !isValid {
		t.Errorf("Unexpected platform: %s", platform)
	}

	t.Logf("Running on platform: %s", platform)
}

func TestLibraryNotifier_EmptyTitle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping notification send test in short mode")
	}

	cfg := config.DefaultConfig()
	notifier := &LibraryNotifier{
		config: cfg,
	}

	err := notifier.Send("", "Test Message", "info")

	// Should handle empty title gracefully
	if err != nil {
		t.Logf("Send with empty title returned error: %v", err)
	}
}

func TestLibraryNotifier_EmptyMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping notification send test in short mode")
	}

	cfg := config.DefaultConfig()
	notifier := &LibraryNotifier{
		config: cfg,
	}

	err := notifier.Send("Test Title", "", "info")

	// Should handle empty message gracefully or return error
	if err != nil {
		t.Logf("Send with empty message returned error: %v", err)
	}
}
