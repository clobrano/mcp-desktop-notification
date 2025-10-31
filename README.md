# MCP Desktop Notification Server (mcp-poke)

A Model Context Protocol (MCP) server that enables AI agents to send desktop notifications to users across Linux, macOS, and Windows.

## Features

- 🔔 **Cross-platform notifications** using [beeep](https://github.com/gen2brain/beeep)
- 🎨 **Severity levels** (info, warning, error, success) with appropriate icons
- ⚙️ **Configurable** via YAML with XDG Base Directory support
- 🧪 **Dry-run mode** for testing without sending actual notifications
- 📝 **Verbose logging** for debugging
- 🔌 **MCP-compatible** using the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk)

## Installation

### From Source

```bash
go install github.com/clobrano/mcp-desktop-notification/cmd/mcp-poke@latest
```

### Building Locally

```bash
git clone https://github.com/clobrano/mcp-desktop-notification
cd mcp-desktop-notification
go build -o mcp-poke ./cmd/mcp-poke
```

## Usage

### As an MCP Server

Add to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "poke": {
      "command": "mcp-poke",
      "args": []
    }
  }
}
```

### Command Line Options

```bash
mcp-poke [options]

Options:
  -config string
        Path to configuration file (default: platform-specific)
  -dry-run
        Dry run mode (log notifications without sending)
  -verbose
        Enable verbose logging
```

### Configuration

The server looks for configuration in platform-specific locations:

- **Linux/macOS**: `~/.config/mcp-desktop-notification/config.yaml`
- **Windows**: `%APPDATA%\mcp-desktop-notification\config.yaml`

See [config.example.yaml](config.example.yaml) for all available options.

## MCP Tool: `poke`

Send a desktop notification to the user.

### Parameters

- `message` (required, string): The notification message text
- `title` (optional, string): The notification title (defaults to "Notification")
- `level` (optional, string): Severity level - one of: `info`, `warning`, `error`, `success` (defaults to "info")

### Example

```json
{
  "name": "poke",
  "arguments": {
    "message": "Your task has completed successfully!",
    "title": "Task Complete",
    "level": "success"
  }
}
```

## Use Cases

- **Long-running tasks**: Notify when data processing, builds, or deployments complete
- **User input required**: Alert when the AI agent needs approval or additional information
- **Error notifications**: Immediately inform users of failures or issues
- **Milestone updates**: Keep users informed of progress in multi-step processes

## Configuration Example

```yaml
notification:
  # Dry-run mode for testing (default: false)
  dry_run: false

  # Verbose logging (default: false, enable only for debugging)
  verbose: false

  # Notification delivery mode (default: "library")
  # Options: "library" (recommended) or "command" (advanced)
  mode: "library"

  # Message template (applies to both modes)
  template:
    default: "{{.Title}}: {{.Message}} [{{.Level}}]"

  # Level mappings for urgency and icons
  levels:
    info:
      urgency: "normal"
      icon: ""
    warning:
      urgency: "normal"
      icon: "dialog-warning"
    error:
      urgency: "critical"
      icon: "dialog-error"
    success:
      urgency: "low"
      icon: "dialog-information"
```

## Development

### Running Tests

```bash
# Run all tests (short mode to skip desktop notification tests)
go test -short ./...

# Run tests with verbose output
go test -short -v ./...

# Test a specific package
go test -short -v ./internal/config
```

### Building

```bash
# Build for current platform
go build -o mcp-poke ./cmd/mcp-poke

# Cross-compile for other platforms
GOOS=linux GOARCH=amd64 go build -o mcp-poke-linux ./cmd/mcp-poke
GOOS=darwin GOARCH=arm64 go build -o mcp-poke-macos ./cmd/mcp-poke
GOOS=windows GOARCH=amd64 go build -o mcp-poke.exe ./cmd/mcp-poke
```

## Requirements

- Go 1.21 or later
- For Linux: D-Bus notification support (usually available by default)
- For macOS: No additional requirements
- For Windows: Windows 10 or later

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Links

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [beeep - Cross-platform notifications](https://github.com/gen2brain/beeep)
