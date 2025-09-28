# Changelog

All notable changes to the Development Tools CLI will be documented in this file.

## [1.0.1] - 2024-01-20

### Added
- Auto-update system with `tools update` command
- Version management and tracking
- Update scripts (PowerShell and Batch)
- Configuration display with `tools config` command
- Environment variable status display

### Changed
- Improved configuration loading from .env files
- Enhanced error handling in all commands

### Fixed
- Windows PATH installation issues
- Environment variable precedence

## [1.0.0] - 2024-01-19

### Added
- Initial release
- PostgreSQL management commands
- MySQL management commands
- RabbitMQ management commands
- MinIO object storage commands
- Environment variable support
- .env file configuration
- Installation scripts for Windows
- Hot reload support with Air

### Features
- Database operations (list, create, drop, backup, restore)
- RabbitMQ queue and exchange management
- MinIO bucket and object operations
- Cross-service configuration management