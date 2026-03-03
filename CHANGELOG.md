# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-03-03

### Added

- Initial release of WarpCTL
- Calendar management commands
  - List calendars
  - Create and delete calendars
  - Manage events with attendees
  - Date range filtering for events
- Contact management commands
  - CRUD operations for contacts
  - Contact group management
  - Bulk operations support
- Task management commands
  - Create and update tasks
  - Priority and due date management
  - Task completion tracking
  - Task folder organization
- Notes management commands
  - Create and edit notes
  - Color-coded organization
  - Tag support for categorization
  - Note folder management
- Spam management commands
  - List quarantined items
  - Deliver spam to inbox
  - Delete spam items
  - Whitelist/Blacklist operations
  - Bulk spam cleanup
- Mail operations via IMAP
  - Clean mailbox command
  - Dry-run mode for safe testing
- System maintenance commands
  - Server authentication
  - Domain management
  - User administration
  - System statistics
- Configuration management
  - YAML-based configuration
  - Environment variable support
  - Config initialization and viewing
- Beautiful table-based output formatting
  - Unicode rounded borders
  - Consistent styling across all commands
  - Professional presentation
- Cross-platform support
  - Linux (AMD64, ARM64)
  - Windows (AMD64)
- Comprehensive documentation
  - Professional README
  - Contributing guidelines
  - MIT License
- CI/CD workflows
  - Automated builds for Linux and Windows
  - Line ending enforcement (LF only)
  - Automated releases on tags
  - Code linting

### Technical Details

- Built with Go 1.21+
- Uses Cobra for CLI framework
- go-pretty for table formatting
- Viper for configuration management
- IMAP support for direct mailbox operations
- XML-RPC for IceWarp API communication

[1.0.0]: https://github.com/rsdenck/warp/releases/tag/v1.0.0
