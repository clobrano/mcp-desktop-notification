# Clarifying Questions for Desktop Notification MCP

Please answer the following questions to help me create a comprehensive PRD. You can reply inline with your selections/answers.

## 1. MCP Server Scope & Integration

**Question:** How should this MCP server be packaged and distributed?

A. As a standalone MCP server that any MCP client can connect to
B. As a library/module that can be embedded into existing MCP servers
C. Both - provide standalone server AND reusable library
D. Other (please specify)

**Your answer:** A

---

## 2. Notification Triggers

**Question:** What specific events should trigger desktop notifications? (Select all that apply)

A. When the AI agent needs user input/approval to proceed
B. When a long-running task completes successfully
C. When a task fails or encounters an error
D. When the agent reaches certain milestones in a task
E. Custom/configurable triggers defined by the user
F. Other (please specify)

**Your answer:** A, B, C, D

---

## 3. Notification Content & Customization

**Question:** What information should be included in notifications?

A. Just a simple message text
B. Message + severity level (info/warning/error/success)
C. Message + severity + action buttons (e.g., "View", "Dismiss", "Retry")
D. Fully customizable with title, body, icon, and actions
E. Other (please specify)

**Your answer:** I want the message to be configurable, as in a Template, with a sane default containing the text passed by the agent, and Info level severity

---

## 4. Configuration Management

**Question:** How should users configure the notification commands and settings?

A. JSON configuration file in a standard location (e.g., `~/.config/mcp-notifications/config.json`)
B. Environment variables
C. Both configuration file AND environment variables (env vars override file)
D. MCP server parameters passed during initialization
E. Other (please specify)

**Your answer:** YAML would be better. Consider using XDG specification for the location

---

## 5. Default Commands per Platform

**Question:** Should the system auto-detect the platform and use sensible defaults?

A. Yes - auto-detect Linux (notify-send), macOS (osascript/terminal-notifier), Windows (PowerShell)
B. Yes - but require explicit configuration before sending any notifications
C. No - always require manual configuration
D. Other (please specify)

**Your answer:** yes

---

## 6. Linux Desktop Environment Detection

**Question:** For Linux, should the system attempt to detect and use the appropriate notification system?

A. Yes - auto-detect KDE (kdialog), GNOME (notify-send), etc.
B. No - just use notify-send as a universal default for Linux
C. Let the user specify their Linux notification command explicitly
D. Other (please specify)

**Your answer:** prefer notify-send if available in the system

---

## 7. Error Handling

**Question:** What should happen if a notification fails to send?

A. Fail silently (log error but don't interrupt the agent)
B. Return an error to the MCP client but continue operation
C. Retry with fallback notification methods
D. Allow configurable behavior (silent/error/retry)
E. Other (please specify)

**Your answer:** B

---

## 8. MCP Tools/Resources

**Question:** What MCP capabilities should this server expose?

A. A single "notify" tool that agents can call
B. Multiple tools for different notification types (notify_info, notify_error, notify_success, etc.)
C. Resources that expose notification history/logs
D. Prompts that help agents decide when to send notifications
E. Combination of the above (please specify which)

**Your answer:** A

---

## 9. Notification Persistence & History

**Question:** Should the system keep a history of sent notifications?

A. Yes - store in a local database/file for auditing
B. Yes - expose via MCP resources so clients can query history
C. No - notifications are ephemeral only
D. Optional/configurable feature
E. Other (please specify)

**Your answer:** No

---

## 10. Testing & Development

**Question:** What testing capabilities are important for development?

A. Dry-run/test mode that shows what would be sent without actually sending
B. Mock notification system for automated testing
C. Verbose logging mode to debug notification delivery
D. All of the above
E. Other (please specify)

**Your answer:** A, we can actually test only on Linux

---

## 11. User Stories Priority

**Question:** Which user stories are most important to you? (Rank 1-5, where 1 is highest priority)

A. ___ As an AI agent, I want to notify the user when I need their input, so they know to check back
B. ___ As an AI agent, I want to notify the user when a long task completes, so they don't have to keep checking
C. ___ As a user, I want to customize notification commands per platform, so I can use my preferred notification system
D. ___ As a user, I want to see notification history, so I can review what the agent has done
E. ___ As a developer, I want to test notifications without spamming my desktop, so I can develop efficiently

**Your answer:** A, B

---

## 12. Additional Requirements

**Question:** Are there any other specific requirements, constraints, or features you want included?

**Your answer:** I want it in Golang

---

## 13. Non-Goals

**Question:** What should this MCP server explicitly NOT do?

A. Handle notification interactions/callbacks (just send, don't process responses)
B. Provide a GUI for configuration
C. Support mobile platforms (iOS/Android)
D. Support non-desktop notification channels (email, SMS, etc.)
E. Other (please specify)

**Your answer:**

---

Please fill in your answers above, and I'll use them to create a comprehensive PRD for your desktop notification MCP server!
