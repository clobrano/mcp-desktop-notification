package mcp

import (
	"testing"

	"github.com/clobrano/mcp-desktop-notification/internal/config"
	"github.com/clobrano/mcp-desktop-notification/internal/notifier"
)

func TestNewServer(t *testing.T) {
	cfg := config.DefaultConfig()
	noti, err := notifier.NewNotifier(cfg)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	server := NewServer(cfg, noti)
	if server == nil {
		t.Error("Expected non-nil server")
	}

	if server.config != cfg {
		t.Error("Server config not set correctly")
	}

	if server.notifier != noti {
		t.Error("Server notifier not set correctly")
	}
}

func TestValidatePokeParams_Valid(t *testing.T) {
	params := map[string]interface{}{
		"message": "Test message",
		"title":   "Test title",
		"level":   "info",
	}

	message, title, level, err := validatePokeParams(params)
	if err != nil {
		t.Errorf("Expected valid params, got error: %v", err)
	}

	if message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", message)
	}

	if title != "Test title" {
		t.Errorf("Expected title 'Test title', got '%s'", title)
	}

	if level != "info" {
		t.Errorf("Expected level 'info', got '%s'", level)
	}
}

func TestValidatePokeParams_MissingMessage(t *testing.T) {
	params := map[string]interface{}{
		"title": "Test title",
	}

	_, _, _, err := validatePokeParams(params)
	if err == nil {
		t.Error("Expected error for missing message parameter")
	}
}

func TestValidatePokeParams_DefaultTitle(t *testing.T) {
	params := map[string]interface{}{
		"message": "Test message",
	}

	message, title, level, err := validatePokeParams(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if title != "Notification" {
		t.Errorf("Expected default title 'Notification', got '%s'", title)
	}

	if level != "info" {
		t.Errorf("Expected default level 'info', got '%s'", level)
	}

	if message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", message)
	}
}

func TestValidatePokeParams_InvalidLevel(t *testing.T) {
	params := map[string]interface{}{
		"message": "Test message",
		"level":   "invalid",
	}

	_, _, _, err := validatePokeParams(params)
	if err == nil {
		t.Error("Expected error for invalid level")
	}
}

func TestValidatePokeParams_EmptyMessage(t *testing.T) {
	params := map[string]interface{}{
		"message": "",
	}

	_, _, _, err := validatePokeParams(params)
	if err == nil {
		t.Error("Expected error for empty message")
	}
}

func TestValidatePokeParams_WrongType(t *testing.T) {
	params := map[string]interface{}{
		"message": 123, // Should be string
	}

	_, _, _, err := validatePokeParams(params)
	if err == nil {
		t.Error("Expected error for wrong parameter type")
	}
}
