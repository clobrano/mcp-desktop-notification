package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Notification.DryRun {
		t.Error("Expected DryRun to be false by default")
	}

	if cfg.Notification.Verbose {
		t.Error("Expected Verbose to be false by default")
	}

	if cfg.Notification.Template.Default == "" {
		t.Error("Expected default template to be set")
	}

	// Check level mappings exist
	if _, ok := cfg.Notification.Levels["info"]; !ok {
		t.Error("Expected 'info' level to be defined")
	}
	if _, ok := cfg.Notification.Levels["warning"]; !ok {
		t.Error("Expected 'warning' level to be defined")
	}
	if _, ok := cfg.Notification.Levels["error"]; !ok {
		t.Error("Expected 'error' level to be defined")
	}
	if _, ok := cfg.Notification.Levels["success"]; !ok {
		t.Error("Expected 'success' level to be defined")
	}
}

func TestValidateConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/path/config.yaml")

	// Should return default config without error when file doesn't exist
	if err != nil {
		t.Errorf("Expected no error for non-existent file, got: %v", err)
	}

	if cfg.Notification.DryRun {
		t.Error("Expected default config when file doesn't exist")
	}
}

func TestLoadConfig_ValidYAML(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
notification:
  dry_run: true
  verbose: true
  template:
    default: "Test: {{.Message}}"
  levels:
    info:
      urgency: normal
      icon: ""
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !cfg.Notification.DryRun {
		t.Error("Expected DryRun to be true")
	}

	if !cfg.Notification.Verbose {
		t.Error("Expected Verbose to be true")
	}

	if cfg.Notification.Template.Default != "Test: {{.Message}}" {
		t.Errorf("Expected custom template, got: %s", cfg.Notification.Template.Default)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `
notification:
  dry_run: true
  invalid: [unclosed
`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestGetConfigPath_XDG(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDG)
		os.Setenv("HOME", originalHome)
	}()

	// Test with XDG_CONFIG_HOME set
	os.Setenv("XDG_CONFIG_HOME", "/custom/config")
	os.Setenv("HOME", "/home/user")

	path := GetConfigPath()

	// Skip on Windows
	if filepath.Separator == '\\' {
		t.Skip("Skipping XDG test on Windows")
	}

	expected := "/custom/config/mcp-desktop-notification/config.yaml"
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestGetConfigPath_XDGFallback(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDG)
		os.Setenv("HOME", originalHome)
	}()

	// Test with XDG_CONFIG_HOME not set
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", "/home/user")

	path := GetConfigPath()

	// Skip on Windows
	if filepath.Separator == '\\' {
		t.Skip("Skipping XDG fallback test on Windows")
	}

	expected := "/home/user/.config/mcp-desktop-notification/config.yaml"
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestGetConfigPath_Windows(t *testing.T) {
	// Skip if not on Windows
	if filepath.Separator != '\\' {
		t.Skip("Skipping Windows test on non-Windows platform")
	}

	// Save original env
	originalAppData := os.Getenv("APPDATA")
	defer os.Setenv("APPDATA", originalAppData)

	os.Setenv("APPDATA", "C:\\Users\\TestUser\\AppData\\Roaming")

	path := GetConfigPath()

	expected := "C:\\Users\\TestUser\\AppData\\Roaming\\mcp-desktop-notification\\config.yaml"
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}
