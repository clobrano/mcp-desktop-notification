# MCP Desktop Notification Server (mcp-poke)

A Model Context Protocol (MCP) server that enables AI agents to send desktop notifications to users across Linux, macOS, and Windows.

## Features

- 🔔 **Cross-platform notifications** using [beeep](https://github.com/gen2brain/beeep)
- 🎨 **Severity levels** (info, warning, error, success) with appropriate icons
- ⚙️ **Configurable** via YAML with XDG Base Directory support
- 🧪 **Dry-run mode** for testing without sending actual notifications
- 📝 **Verbose logging** for debugging
- 🔌 **MCP-compatible** using the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- 📂 **Workspace identification** - app name displays the last 2 directories from PWD

## Installation

### From Source

```bash
go install github.com/clobrano/mcp-desktop-notification@latest
```

### Building Locally

```bash
git clone https://github.com/clobrano/mcp-desktop-notification
cd mcp-desktop-notification
go build -o mcp-poke .
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

### Examples

**Simple notification:**
```json
{
  "name": "poke",
  "arguments": {
    "message": "Your task has completed successfully!"
  }
}
```
Result: Desktop notification with title "Notification" and your message.

**Task completion with custom title:**
```json
{
  "name": "poke",
  "arguments": {
    "message": "Data processing finished. 10,000 records processed.",
    "title": "Task Complete",
    "level": "success"
  }
}
```
Result: Green/success notification with custom title.

**Warning notification:**
```json
{
  "name": "poke",
  "arguments": {
    "message": "Database migration requires your approval to continue.",
    "title": "Action Required",
    "level": "warning"
  }
}
```
Result: Yellow/warning notification.

**Error alert:**
```json
{
  "name": "poke",
  "arguments": {
    "message": "Failed to connect to API. Check your network connection.",
    "title": "Connection Error",
    "level": "error"
  }
}
```
Result: Red/critical notification.

## Use Cases

- **Long-running tasks**: Notify when data processing, builds, or deployments complete
- **User input required**: Alert when the AI agent needs approval or additional information
- **Error notifications**: Immediately inform users of failures or issues
- **Milestone updates**: Keep users informed of progress in multi-step processes

## Workspace Identification

The notification app name automatically displays the last 2 directories from your `PWD` (current working directory) to help you identify which workspace is sending notifications.

**Examples:**
- Working in `/home/carlo/workspace/foo` → Notifications show app name as `workspace/foo`
- Working in `/home/carlo/projects/my-app` → Notifications show app name as `projects/my-app`
- Working in `/projects` → Notifications show app name as `projects`

This is especially useful when running multiple AI agents in different workspaces, allowing you to quickly identify the source of each notification.

**Default:** If PWD is not available or is the root directory, the app name defaults to `mcp-poke`.

## Configuration

### Message Templates

**Note**: The current implementation uses the message and title directly as provided by the agent. Template support is planned for a future release.

The configuration file supports customization of notification behavior:

```yaml
notification:
  # Dry-run mode for testing (default: false)
  dry_run: false

  # Verbose logging (default: false, enable only for debugging)
  verbose: false

  # Notification delivery mode (default: "library")
  # Options: "library" (recommended) or "command" (advanced)
  mode: "library"

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

### Example: Customizing Notification Levels

```yaml
notification:
  levels:
    error:
      urgency: "critical"       # Makes error notifications more prominent
      icon: "dialog-error"       # Uses red error icon
    success:
      urgency: "low"             # Success notifications are less intrusive
      icon: "emblem-default"     # Uses a checkmark or success icon
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
go build -o mcp-poke .

# Cross-compile for other platforms
GOOS=linux GOARCH=amd64 go build -o mcp-poke-linux .
GOOS=darwin GOARCH=arm64 go build -o mcp-poke-macos .
GOOS=windows GOARCH=amd64 go build -o mcp-poke.exe .
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
