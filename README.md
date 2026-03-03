# WarpCTL

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)
![Go Version](https://img.shields.io/badge/go-1.21+-blue?style=flat-square)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20windows-blue?style=flat-square)

A powerful command-line interface for managing IceWarp server operations.

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Documentation](#documentation)

</div>

---

## Overview

WarpCTL is a comprehensive CLI tool designed for efficient management of IceWarp mail server infrastructure. Built with Go, it provides a unified interface for calendar, contacts, tasks, notes, spam management, and system maintenance operations.

## Features

### Core Capabilities

<table>
<tr>
<td width="50%">

**Calendar Management**
- List and manage calendars
- Create and delete events
- Schedule management
- Attendee coordination

**Contact Management**
- Contact CRUD operations
- Group management
- Bulk operations support
- Contact folder organization

</td>
<td width="50%">

**Task Management**
- Task creation and tracking
- Priority management
- Due date handling
- Task folder organization

**Notes Management**
- Note creation and editing
- Color-coded organization
- Tag support
- Folder management

</td>
</tr>
<tr>
<td width="50%">

**Spam Control**
- Quarantine management
- Whitelist/Blacklist operations
- Bulk spam cleanup
- Message delivery control

</td>
<td width="50%">

**System Maintenance**
- Server authentication
- Domain management
- User administration
- System monitoring

</td>
</tr>
</table>

### Technical Features

- **Formatted Output**: Beautiful table-based output using Unicode borders
- **IMAP Support**: Direct mailbox operations via IMAP protocol
- **Configuration Management**: Flexible YAML-based configuration
- **Environment Variables**: Support for credential management via env vars
- **Cross-Platform**: Native support for Linux and Windows

## Installation

### Prerequisites

- Go 1.21 or higher
- IceWarp server access credentials
- Network connectivity to IceWarp server

### From Source

```bash
git clone https://github.com/rsdenck/warp.git
cd warp
go build -o warpctl
```

### Binary Installation

Download the latest release for your platform from the [releases page](https://github.com/rsdenck/warp/releases).

#### Linux

```bash
chmod +x warpctl
sudo mv warpctl /usr/local/bin/
```

#### Windows

Add the executable to your PATH or run directly from the download location.

## Configuration

### Initialize Configuration

```bash
warpctl config init
```

### Set Server Details

```bash
warpctl config set server.url https://your-icewarp-server.com
warpctl config set auth.username your-email@domain.com
warpctl config set auth.password your-password
```

### Environment Variables

Alternatively, use environment variables:

```bash
export IW_USERNAME="your-email@domain.com"
export IW_PASSWORD="your-password"
```

### Configuration File Location

- Linux: `~/.icwli/icwli.yaml`
- Windows: `%USERPROFILE%/.icwli/icwli.yaml`

## Usage

### Calendar Operations

```bash
# Login to calendar API
warpctl calendar login

# List calendars
warpctl calendar list

# Create a calendar
warpctl calendar create "Work Calendar" --description "Work events" --color "#FF0000"

# List events
warpctl calendar events <calendar-id> --start 2024-01-01 --end 2024-12-31

# Create an event
warpctl calendar create-event <calendar-id> "Meeting" \
  --start "2024-03-15T14:00" \
  --end "2024-03-15T15:00" \
  --location "Conference Room" \
  --attendee "user@domain.com"
```

### Contact Management

```bash
# List contacts
warpctl contacts list

# Create a contact
warpctl contacts create user@example.com \
  --first-name "John" \
  --last-name "Doe" \
  --phone "+1234567890" \
  --company "Acme Corp"

# List contact groups
warpctl contacts groups

# Create a group
warpctl contacts create-group "Sales Team"
```

### Task Management

```bash
# List tasks
warpctl tasks list

# Create a task
warpctl tasks create "Complete project" \
  --description "Finish the Q1 project" \
  --priority 1 \
  --due "2024-03-31"

# Update a task
warpctl tasks update <task-id> \
  --status "in-progress" \
  --priority 2

# Mark task as complete
warpctl tasks complete <task-id>
```

### Notes Management

```bash
# List notes
warpctl notes list

# Create a note
warpctl notes create "Meeting Notes" \
  --content "Discussion points..." \
  --color "#FFFF00" \
  --tag "work" \
  --tag "meeting"

# Update a note
warpctl notes update <note-id> \
  --title "Updated Title" \
  --content "New content"

# List note folders
warpctl notes folders
```

### Spam Management

```bash
# Login to spam API
warpctl spam login

# List quarantined items
warpctl spam list --limit 50

# Get spam item details
warpctl spam info <item-id>

# Deliver spam to inbox
warpctl spam deliver <item-id>

# Delete spam item
warpctl spam delete <item-id>

# Clean all spam
warpctl spam clean

# Whitelist a sender
warpctl spam whitelist sender@example.com

# Blacklist a sender
warpctl spam blacklist sender@example.com
```

### Mailbox Operations

```bash
# Clean a mailbox (via IMAP)
warpctl clean --mailbox "Spam" --yes

# Dry run (preview only)
warpctl clean --mailbox "Spam" --dry-run
```

### System Maintenance

```bash
# Login to maintenance API
warpctl maintenance login

# Get server information
warpctl maintenance info

# List domains
warpctl maintenance domains

# List users in a domain
warpctl maintenance users <domain>

# Get system statistics
warpctl maintenance stats
```

## Output Format

All list commands display data in formatted tables with Unicode borders:

```
╭──────────────────────────────────────────────────╮
│                    CONTACTS                      │
├─────────────┬─────────────┬──────────────┬───────┤
│ First Name  │ Last Name   │ Email        │ ID    │
│ John        │ Doe         │ john@ex.com  │ 12345 │
│ Jane        │ Smith       │ jane@ex.com  │ 12346 │
╰─────────────┴─────────────┴──────────────┴───────╯
```

## Development

### Project Structure

```
warpctl/
├── cmd/                    # Command implementations
│   ├── calendar/          # Calendar commands
│   ├── contacts/          # Contact commands
│   ├── tasks/             # Task commands
│   ├── notes/             # Notes commands
│   ├── spam/              # Spam commands
│   ├── mail/              # Mail commands
│   ├── clean/             # Cleanup commands
│   ├── maintenance/       # Maintenance commands
│   └── root/              # Root command
├── internal/              # Internal packages
│   ├── config/           # Configuration management
│   ├── imap/             # IMAP client
│   ├── logger/           # Logging utilities
│   ├── output/           # Output formatting
│   └── sdk/              # IceWarp SDK
├── .github/              # GitHub workflows
├── go.mod                # Go module definition
├── go.sum                # Go dependencies
├── LICENSE               # MIT License
├── README.md             # This file
└── main.go               # Application entry point
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rsdenck/warp.git
cd warp

# Install dependencies
go mod download

# Build for current platform
go build -o warpctl

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o warpctl-linux-amd64
GOOS=windows GOARCH=amd64 go build -o warpctl-windows-amd64.exe
```

### Running Tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices and idioms
- Use `gofmt` for code formatting
- Add tests for new features
- Update documentation as needed

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:

1. Verify your credentials are correct
2. Ensure the server URL is accessible
3. Check if your account has necessary permissions
4. Try using environment variables instead of config file

### Connection Issues

```bash
# Test server connectivity
curl -I https://your-icewarp-server.com

# Enable debug mode
warpctl --debug <command>
```

### Configuration Issues

```bash
# View current configuration
warpctl config view

# Reinitialize configuration
rm ~/.icwli/icwli.yaml
warpctl config init
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Table formatting by [go-pretty](https://github.com/jedib0t/go-pretty)
- Configuration management by [Viper](https://github.com/spf13/viper)

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/rsdenck/warp).

---

<div align="center">

Made with Go | Maintained by [Ranlens Denck](https://github.com/rsdenck)

</div>
