package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Notification NotificationConfig `yaml:"notification"`
}

// NotificationConfig contains notification-specific settings
type NotificationConfig struct {
	DryRun   bool              `yaml:"dry_run"`
	Verbose  bool              `yaml:"verbose"`
	Template Template          `yaml:"template"`
	Levels   map[string]Level  `yaml:"levels"`
}

// Template contains message template configuration
type Template struct {
	Default string `yaml:"default"`
}

// Level contains configuration for a notification severity level
type Level struct {
	Urgency string `yaml:"urgency"`
	Icon    string `yaml:"icon"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Notification: NotificationConfig{
			DryRun:  false,
			Verbose: false,
			Template: Template{
				Default: "{{.Title}}: {{.Message}} [{{.Level}}]",
			},
			Levels: map[string]Level{
				"info": {
					Urgency: "normal",
					Icon:    "",
				},
				"warning": {
					Urgency: "normal",
					Icon:    "dialog-warning",
				},
				"error": {
					Urgency: "critical",
					Icon:    "dialog-error",
				},
				"success": {
					Urgency: "low",
					Icon:    "dialog-information",
				},
			},
		},
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// No validation needed for now
	// Future: could validate level names, icon paths, etc.
	return nil
}

// LoadConfig loads configuration from a file, or returns defaults if file doesn't exist
func LoadConfig(path string) (*Config, error) {
	// If file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Start with default config
	cfg := DefaultConfig()

	// Unmarshal YAML, merging with defaults
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate the config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// GetConfigPath returns the platform-specific config file path
func GetConfigPath() string {
	if runtime.GOOS == "windows" {
		// Windows: %APPDATA%\mcp-desktop-notification\config.yaml
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		return filepath.Join(appData, "mcp-desktop-notification", "config.yaml")
	}

	// Linux/macOS: XDG Base Directory specification
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(xdgConfigHome, "mcp-desktop-notification", "config.yaml")
}

// LoadDefaultConfig loads config from the default platform-specific path
func LoadDefaultConfig() (*Config, error) {
	return LoadConfig(GetConfigPath())
}
