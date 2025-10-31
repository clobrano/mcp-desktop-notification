# Product Requirements Document: Desktop Notification MCP Server

## 1. Introduction/Overview

This document outlines the requirements for a Model Context Protocol (MCP) server that enables AI agents to send desktop notifications to users across multiple operating systems (Linux, macOS, and Windows). The server will notify users when agents need input, when tasks complete, when errors occur, or when significant milestones are reached.

**Problem Statement:** AI agents often perform long-running tasks or encounter situations requiring user input, but users have no way of knowing when to check back without actively monitoring the terminal or application. This leads to wasted time and poor user experience.

**Solution:** A cross-platform MCP server that allows agents to send desktop notifications, keeping users informed of important events without requiring constant monitoring.

## 2. Goals

1. Provide a standalone MCP server that any MCP client can connect to for sending desktop notifications
2. Support Linux (KDE/GNOME), macOS, and Windows with platform-appropriate notification systems
3. Allow users to customize notification commands and message templates via configuration
4. Deliver notifications reliably with clear error reporting when delivery fails
5. Implement the solution in Go for cross-platform compatibility and easy distribution

## 3. User Stories

### Priority 1 (Highest)
- **US-1:** As an AI agent, I want to notify the user when I need their input, so they know to check back and provide the necessary information
- **US-2:** As an AI agent, I want to notify the user when a long task completes, so they don't have to keep checking the terminal

### Priority 2
- **US-3:** As an AI agent, I want to notify the user when a task fails or encounters an error, so they can intervene quickly
- **US-4:** As an AI agent, I want to notify the user when I reach certain milestones in a task, so they can track progress
- **US-5:** As a user, I want to customize notification commands per platform, so I can use my preferred notification system
- **US-6:** As a user, I want to customize notification message templates, so notifications match my preferences
- **US-7:** As a developer, I want to test notifications without spamming my desktop, so I can develop efficiently

## 4. Functional Requirements

### 4.1 MCP Server Setup
- **FR-1:** The system MUST be packaged as a standalone MCP server that can be launched and connected to by any MCP client
- **FR-2:** The system MUST be implemented in Go (Golang)
- **FR-3:** The system MUST provide a single executable binary for each supported platform

### 4.2 Platform Support & Notification Delivery
- **FR-4:** The system MUST support Linux, macOS, and Windows operating systems
- **FR-5:** The system MUST auto-detect the current platform at runtime
- **FR-6:** The system SHOULD use a Go notification library (e.g., `beeep`, `notify`) as the primary notification mechanism to provide cross-platform abstraction without external dependencies
- **FR-7:** If Go libraries are insufficient or unavailable, the system MAY fall back to platform-specific command execution
- **FR-8:** Users MUST be able to override the default notification mechanism with custom commands via configuration (see FR-15)

### 4.3 MCP Tool Interface
- **FR-9:** The system MUST expose a single MCP tool named `poke` that agents can call
- **FR-10:** The `poke` tool MUST accept the following parameters:
  - `message` (required, string): The notification message text
  - `title` (optional, string): The notification title
  - `level` (optional, string): Severity level (info/warning/error/success), defaults to "info"
- **FR-11:** The `poke` tool MUST return success confirmation or error details

### 4.4 Configuration Management
- **FR-12:** The system MUST support configuration via YAML file
- **FR-13:** The configuration file MUST be located following XDG Base Directory specification:
  - Linux/macOS: `$XDG_CONFIG_HOME/mcp-desktop-notification/config.yaml` (fallback to `~/.config/mcp-desktop-notification/config.yaml`)
  - Windows: `%APPDATA%\mcp-desktop-notification\config.yaml`
- **FR-14:** The system MUST provide sensible default behavior if no config file exists (using Go library notifications)
- **FR-15:** The configuration MUST allow users to optionally specify custom notification commands per platform to override the default Go library behavior
- **FR-16:** The configuration MUST support customizable message templates with variable substitution (applicable to both library-based and command-based notifications)

### 4.5 Message Template System
- **FR-17:** The system MUST support message templates with the following variables:
  - `{{.Message}}`: The message text passed by the agent
  - `{{.Title}}`: The title passed by the agent
  - `{{.Level}}`: The severity level
  - `{{.Timestamp}}`: The notification timestamp
- **FR-18:** The system MUST provide a sane default template that includes the message text and severity level
- **FR-19:** Users MUST be able to override the default template in the configuration file

### 4.6 Notification Triggers
- **FR-20:** The system MUST support notifications for the following scenarios:
  - When the AI agent needs user input/approval to proceed
  - When a long-running task completes successfully
  - When a task fails or encounters an error
  - When the agent reaches certain milestones in a task

### 4.7 Error Handling
- **FR-21:** If a notification fails to send, the system MUST return an error to the MCP client
- **FR-22:** The system MUST continue operation after a notification failure (non-blocking)
- **FR-23:** Error messages MUST include clear information about why the notification failed
- **FR-24:** The system MUST log all notification attempts and their outcomes

### 4.8 Testing & Development
- **FR-25:** The system MUST provide a dry-run/test mode that shows what would be sent without actually sending notifications
- **FR-26:** Dry-run mode MUST be configurable via command-line flag or configuration file
- **FR-27:** The system MUST provide verbose logging mode to debug notification delivery
- **FR-28:** Verbose logging MUST be disabled by default and only enabled explicitly for testing purposes via command-line flag or configuration file

## 5. Non-Goals (Out of Scope)

The following features are explicitly NOT included in this version:

- **NG-1:** Handling notification interactions/callbacks (the system will only send notifications, not process user responses to them)
- **NG-2:** Providing a GUI for configuration (configuration is file-based only)
- **NG-3:** Supporting mobile platforms (iOS/Android)
- **NG-4:** Supporting non-desktop notification channels (email, SMS, Slack, etc.)
- **NG-5:** Maintaining notification history or logs accessible via MCP resources
- **NG-6:** Persistent storage or database integration
- **NG-7:** Multiple notification tools (only one `poke` tool will be provided)
- **NG-8:** Full testing support on macOS and Windows (testing focus is Linux only)

## 6. Design Considerations

### 6.1 Configuration File Structure (YAML)

```yaml
# Example configuration file structure
notification:
  # Dry-run mode for testing (default: false)
  dry_run: false

  # Verbose logging for debugging (default: false, only enable for testing)
  verbose: false

  # Notification delivery mode
  # Options: "library" (default, uses Go notification library) or "command" (uses custom commands)
  mode: "library"

  # Optional: Platform-specific custom commands (only used if mode is "command")
  # If not specified, the system will use Go notification libraries (e.g., beeep)
  commands:
    linux: "notify-send '{{.Title}}' '{{.Message}}' -u {{.UrgencyLevel}}'"
    macos: "osascript -e 'display notification \"{{.Message}}\" with title \"{{.Title}}\"'"
    windows: "powershell -Command \"[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null; ...\""

  # Message template (applies to both library and command modes)
  template:
    default: "{{.Title}}: {{.Message}} [{{.Level}}]"

  # Level mapping (for urgency, icons, etc. - used by library mode and available for command templates)
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

### 6.2 MCP Tool Schema

```json
{
  "name": "poke",
  "description": "Send a desktop notification to the user",
  "inputSchema": {
    "type": "object",
    "properties": {
      "message": {
        "type": "string",
        "description": "The notification message text"
      },
      "title": {
        "type": "string",
        "description": "The notification title (optional)"
      },
      "level": {
        "type": "string",
        "enum": ["info", "warning", "error", "success"],
        "description": "Severity level (optional, defaults to 'info')"
      }
    },
    "required": ["message"]
  }
}
```

## 7. Technical Considerations

### 7.1 Technology Stack
- **Language:** Go (Golang) for cross-platform support and easy distribution
- **MCP SDK:** Use official MCP SDK for Go (or implement MCP protocol directly if SDK unavailable)
- **Configuration:** Use Go YAML library (e.g., `gopkg.in/yaml.v3`)
- **Template Engine:** Use Go's `text/template` for message templating
- **Notification Libraries:** Prefer Go native notification libraries if available (e.g., `github.com/gen2brain/beeep`, `github.com/martinlindhe/notify`) as they provide cross-platform abstractions and eliminate the need for external command execution

### 7.2 Platform Detection & Notification Delivery
- **Primary approach:** Use Go notification libraries (e.g., `beeep`) that provide cross-platform abstractions
  - No external dependencies required (notify-send, osascript, etc.)
  - Built-in support for all platforms
  - No command injection risk
- **Optional approach:** Allow users to configure custom command execution if they need specific features not provided by Go libraries
- Use Go's `runtime.GOOS` to detect the operating system for any platform-specific behavior
- Evaluate Go libraries for security, feature completeness, active maintenance, and cross-platform reliability

### 7.3 Dependencies
- Minimize external dependencies for easier distribution
- Evaluate Go notification library dependency footprint during library selection
- **Advantage of Go libraries:** No external runtime dependencies (no need for notify-send, osascript, PowerShell scripts, etc.)
- If users configure custom commands, those commands become their responsibility to install/maintain

### 7.4 Security
- **Go library approach:** Inherently safe from command injection since no shell commands are executed
- **Custom command approach:** If users configure custom commands:
  - Sanitize user input to prevent command injection
  - Validate template variables to prevent template injection attacks
  - Implement proper shell escaping for command arguments
  - Document security considerations for custom command configuration

## 8. Success Metrics

- **SM-1:** Successful notification delivery rate > 95% using Go notification libraries (no external dependencies required)
- **SM-2:** Cross-platform compatibility verified on Linux (Ubuntu/Fedora with GNOME/KDE), macOS, and Windows without requiring any system-specific notification tools
- **SM-3:** Average notification latency < 500ms from MCP tool call to desktop notification
- **SM-4:** Zero security vulnerabilities (Go library approach eliminates command injection risk)
- **SM-5:** Configuration system works correctly across all platforms following XDG/platform conventions
- **SM-6:** Custom command configuration works correctly when users choose to override defaults

## 9. Open Questions

1. **Q1:** Should the server support checking if notification capabilities are available before attempting to send? (e.g., a "check_availability" tool)
2. **Q2:** Should we support notification icons/images beyond the default level-based icons provided by Go libraries?
3. **Q3:** Which Go notification library should we use as primary implementation? (`beeep` appears most popular, but requires evaluation)
4. **Q4:** Should we provide a command-line interface to send test notifications directly (outside of MCP)?
5. **Q5:** What logging library should be used, or should we use Go's standard `log` package?
6. **Q6:** What should be the default behavior for the `mode` configuration? Always default to "library" unless custom commands are configured?
7. **Q7:** Should template variables work differently between library mode (limited variables) vs command mode (full template flexibility)?

## 10. Implementation Phases (Suggested)

### Phase 1: MVP (Minimum Viable Product)
- Basic MCP server setup in Go
- Single `poke` tool implementation
- Integrate Go notification library (e.g., `beeep`) for cross-platform support
- Basic configuration file loading (YAML) with XDG support
- Simple message passing with title, message, and level parameters

### Phase 2: Advanced Features
- Message template system
- Level-based urgency/icon mapping (as supported by chosen Go library)
- Dry-run/test mode
- Verbose logging (disabled by default, enabled via config/flag for testing)
- Error handling and reporting

### Phase 3: Custom Command Support
- Optional custom command execution mode
- Platform-specific command configuration
- Template variable expansion for custom commands
- Command injection prevention and security measures

### Phase 4: Polish & Distribution
- Comprehensive testing on all platforms (Linux with GNOME/KDE, macOS, Windows)
- Documentation (README, configuration examples, security guidelines)
- Binary distribution setup (GitHub Releases, package managers)
- Performance optimization

## Appendix A: Example Use Cases

### Use Case 1: Agent Needs User Input
```
Agent calls: poke(message="I need your approval to proceed with database migration", title="User Input Required", level="warning")
User sees: Desktop notification prompting them to check the terminal
```

### Use Case 2: Long Task Completion
```
Agent calls: poke(message="Data processing completed successfully. 10,000 records processed.", title="Task Complete", level="success")
User sees: Desktop notification informing them the task is done
```

### Use Case 3: Error Notification
```
Agent calls: poke(message="Failed to connect to API endpoint. Check your network connection.", title="Error", level="error")
User sees: Desktop notification with error icon/urgency
```

## Appendix B: Platform-Specific Notes

### Using Go Notification Libraries (Recommended Default)

Go notification libraries (e.g., `beeep`) handle platform-specific implementation details automatically:

**Linux:**
- Library handles D-Bus notifications natively
- Works across GNOME, KDE, XFCE, and other desktop environments
- No need for notify-send or other external tools

**macOS:**
- Library uses native macOS notification APIs
- Notifications appear in Notification Center
- No need for osascript or terminal-notifier

**Windows:**
- Library uses Windows Toast notification APIs
- Works on Windows 10 and later
- No need for PowerShell scripts

### Using Custom Commands (Optional Override)

If users configure custom commands for specific requirements:

**Linux:**
- `notify-send` works across most desktop environments (GNOME, KDE, XFCE, etc.)
- Supports urgency levels: low, normal, critical
- Supports custom icons
- Must be installed separately

**macOS:**
- `osascript` is built-in and requires no additional dependencies
- `terminal-notifier` provides more features but requires installation
- Notifications appear in Notification Center

**Windows:**
- PowerShell can trigger Windows Toast notifications (Windows 10+)
- Requires appropriate PowerShell commands to create notification objects
- May need to handle Windows version differences

## Document Version

- **Version:** 1.3
- **Last Updated:** 2025-10-31
- **Status:** Draft - Ready for Implementation
- **Changelog:**
  - v1.3: Clarified that verbose logging is disabled by default and only enabled for testing
  - v1.2: Renamed MCP tool from `notify` to `poke` to avoid naming conflicts
  - v1.1: Updated functional requirements to prioritize Go notification libraries over command execution
  - v1.0: Initial draft
