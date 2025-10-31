package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/clobrano/mcp-desktop-notification/internal/config"
	"github.com/clobrano/mcp-desktop-notification/internal/notifier"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP server for desktop notifications
type Server struct {
	config   *config.Config
	notifier notifier.Notifier
	mcp      *mcp.Server
}

// PokeArgs represents the arguments for the poke tool
type PokeArgs struct {
	Message string `json:"message" jsonschema:"The notification message text"`
	Title   string `json:"title,omitempty" jsonschema:"The notification title"`
	Level   string `json:"level,omitempty" jsonschema:"Severity level: info, warning, error, or success"`
}

// NewServer creates a new MCP server
func NewServer(cfg *config.Config, noti notifier.Notifier) *Server {
	return &Server{
		config:   cfg,
		notifier: noti,
	}
}

// Start initializes and starts the MCP server
func (s *Server) Start() error {
	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-poke",
		Version: "1.0.0",
	}, nil)

	s.mcp = mcpServer

	// Register the poke tool
	s.registerPokeToolHandler()

	if s.config.Notification.Verbose {
		log.Println("[MCP Server] Starting MCP server on stdio")
	}

	// Start the server (blocking)
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// registerPokeToolHandler registers the poke tool with the MCP server
func (s *Server) registerPokeToolHandler() {
	// Define the poke tool using AddTool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "poke",
		Description: "Send a desktop notification to communicate with the user. Use for task completions, errors, warnings, or whenever you need the user's attention while they may be working in another application.",
	}, s.handlePokeTool)
}

// handlePokeTool handles the poke tool invocation
func (s *Server) handlePokeTool(ctx context.Context, req *mcp.CallToolRequest, args PokeArgs) (*mcp.CallToolResult, any, error) {
	// Validate and extract parameters
	message, title, level, err := validatePokeArgs(args)
	if err != nil {
		if s.config.Notification.Verbose {
			log.Printf("[MCP Server] Parameter validation error: %v", err)
		}
		// Return error
		return nil, nil, err
	}

	// Log the request if verbose
	if s.config.Notification.Verbose {
		log.Printf("[MCP Server] Received poke request - Title: %s, Message: %s, Level: %s", title, message, level)
	}

	// Send notification
	if err := s.notifier.Send(title, message, level); err != nil {
		errMsg := fmt.Sprintf("Failed to send notification: %v", err)
		if s.config.Notification.Verbose {
			log.Printf("[MCP Server] %s", errMsg)
		}
		return nil, nil, fmt.Errorf("%s", errMsg)
	}

	// Return success
	if s.config.Notification.Verbose {
		log.Printf("[MCP Server] Notification sent successfully")
	}

	successMsg := fmt.Sprintf("Notification sent: %s - %s [%s]", title, message, level)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: successMsg},
		},
	}, nil, nil
}

// validatePokeArgs validates and extracts parameters from PokeArgs
func validatePokeArgs(args PokeArgs) (message, title, level string, err error) {
	// Validate message (required)
	if args.Message == "" {
		return "", "", "", fmt.Errorf("message cannot be empty")
	}
	message = args.Message

	// Set title (optional, default to "Notification")
	title = args.Title
	if title == "" {
		title = "Notification"
	}

	// Set level (optional, default to "info")
	level = args.Level
	if level == "" {
		level = "info"
	}

	// Validate level
	validLevels := map[string]bool{
		"info":    true,
		"warning": true,
		"error":   true,
		"success": true,
	}
	if !validLevels[level] {
		return "", "", "", fmt.Errorf("invalid level: %s (must be one of: info, warning, error, success)", level)
	}

	return message, title, level, nil
}

// validatePokeParams validates and extracts parameters from the poke tool call (for backward compatibility with tests)
func validatePokeParams(params map[string]interface{}) (message, title, level string, err error) {
	// Extract message (required)
	msgVal, ok := params["message"]
	if !ok {
		return "", "", "", fmt.Errorf("missing required parameter: message")
	}

	message, ok = msgVal.(string)
	if !ok {
		return "", "", "", fmt.Errorf("message must be a string")
	}

	if message == "" {
		return "", "", "", fmt.Errorf("message cannot be empty")
	}

	// Extract title (optional, default to "Notification")
	title = "Notification"
	if titleVal, ok := params["title"]; ok {
		if titleStr, ok := titleVal.(string); ok {
			title = titleStr
		}
	}

	// Extract level (optional, default to "info")
	level = "info"
	if levelVal, ok := params["level"]; ok {
		if levelStr, ok := levelVal.(string); ok {
			// Validate level
			validLevels := map[string]bool{
				"info":    true,
				"warning": true,
				"error":   true,
				"success": true,
			}
			if !validLevels[levelStr] {
				return "", "", "", fmt.Errorf("invalid level: %s (must be one of: info, warning, error, success)", levelStr)
			}
			level = levelStr
		}
	}

	return message, title, level, nil
}
