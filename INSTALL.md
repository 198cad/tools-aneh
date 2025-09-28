# Installation Guide - Development Tools CLI

## Quick Installation

### Method 1: Using PowerShell (Recommended)

1. Open PowerShell as Administrator
2. Navigate to the tools directory:
   ```powershell
   cd C:\Users\mafiaboy\Projects\tools
   ```
3. Run the installation script:
   ```powershell
   Set-ExecutionPolicy Bypass -Scope Process -Force
   .\install.ps1
   ```

### Method 2: Using Command Prompt

1. Open Command Prompt as Administrator
2. Navigate to the tools directory:
   ```cmd
   cd C:\Users\mafiaboy\Projects\tools
   ```
3. Run the installation script:
   ```cmd
   install.bat
   ```

## Manual Installation

If the scripts don't work, you can install manually:

### Step 1: Build the executable
```bash
go build -o tools.exe .
```

### Step 2: Add to Windows PATH

#### Option A: Using System Properties GUI
1. Press `Win + X` and select "System"
2. Click "Advanced system settings"
3. Click "Environment Variables"
4. Under "System variables", find and select "Path", then click "Edit"
5. Click "New" and add: `C:\Users\mafiaboy\Projects\tools`
6. Click "OK" on all windows

#### Option B: Using PowerShell (Administrator)
```powershell
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$newPath = $currentPath + ";C:\Users\mafiaboy\Projects\tools"
[Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
```

#### Option C: Using Command Prompt (Administrator)
```cmd
setx /M PATH "%PATH%;C:\Users\mafiaboy\Projects\tools"
```

### Step 3: Configure environment variables
1. Copy `.env.example` to `.env`:
   ```cmd
   copy .env.example .env
   ```
2. Edit `.env` with your configuration

## Verify Installation

After installation, restart your terminal and test:

```cmd
# Check version
tools --version

# Check configuration
tools config all

# Test database commands
tools db list

# Test RabbitMQ commands
tools rabbit queues

# Test MinIO commands
tools minio buckets
```

## Global Environment Variables Setup

To make configuration available globally:

### Option 1: Set System Environment Variables

Using PowerShell (Administrator):
```powershell
# PostgreSQL
[Environment]::SetEnvironmentVariable("PGHOST", "localhost", "Machine")
[Environment]::SetEnvironmentVariable("PGPORT", "5432", "Machine")
[Environment]::SetEnvironmentVariable("PGUSER", "postgres", "Machine")
[Environment]::SetEnvironmentVariable("PGPASSWORD", "yourpassword", "Machine")



# RabbitMQ
[Environment]::SetEnvironmentVariable("RABBITMQ_HOST", "localhost", "Machine")
[Environment]::SetEnvironmentVariable("RABBITMQ_MANAGEMENT_PORT", "15672", "Machine")
[Environment]::SetEnvironmentVariable("RABBITMQ_DEFAULT_USER", "guest", "Machine")
[Environment]::SetEnvironmentVariable("RABBITMQ_DEFAULT_PASS", "guest", "Machine")

# MinIO
[Environment]::SetEnvironmentVariable("MINIO_ENDPOINT", "localhost:9000", "Machine")
[Environment]::SetEnvironmentVariable("MINIO_ACCESS_KEY", "minioadmin", "Machine")
[Environment]::SetEnvironmentVariable("MINIO_SECRET_KEY", "minioadmin", "Machine")
```

### Option 2: User Profile Configuration

Add to your PowerShell profile (`$PROFILE`):
```powershell
# PostgreSQL
$env:PGHOST = "localhost"
$env:PGPORT = "5432"
$env:PGUSER = "postgres"
$env:PGPASSWORD = "yourpassword"



# RabbitMQ
$env:RABBITMQ_HOST = "localhost"
$env:RABBITMQ_MANAGEMENT_PORT = "15672"
$env:RABBITMQ_DEFAULT_USER = "guest"
$env:RABBITMQ_DEFAULT_PASS = "guest"

# MinIO
$env:MINIO_ENDPOINT = "localhost:9000"
$env:MINIO_ACCESS_KEY = "minioadmin"
$env:MINIO_SECRET_KEY = "minioadmin"
```

## Uninstallation

To remove the tools from PATH:

### Using the uninstall script:
```cmd
uninstall.bat
```

### Or manually:
1. Remove `C:\Users\mafiaboy\Projects\tools` from PATH environment variable
2. Delete the folder if desired

## Updating Tools

The tools include a robust auto-update system to keep your installation current.

### Manual Update

#### Method 1: Using the built-in command
```bash
# Check for updates
tools update --check

# Install updates
tools update

# Force update even if already up to date
tools update --force
```

#### Method 2: Using update scripts
```powershell
# PowerShell
.\update.ps1

# Command Prompt
update.bat
```

### Automatic Updates

#### Set up scheduled update checks:

1. **Using Task Scheduler:**
   - Open Task Scheduler
   - Create Basic Task → "Check Tools Updates"
   - Trigger: Daily or Weekly
   - Action: Start `C:\Users\mafiaboy\Projects\tools\auto-update.bat`

2. **Using PowerShell scheduled job:**
```powershell
$trigger = New-JobTrigger -Daily -At "09:00"
$action = {
    Set-Location "C:\Users\mafiaboy\Projects\tools"
    & .\tools.exe update --check
}
Register-ScheduledJob -Name "CheckToolsUpdate" -Trigger $trigger -ScriptBlock $action
```

### Update Features

- **Safe updates**: Backs up configuration and current version
- **Conflict resolution**: Automatically handles merge conflicts
- **Configuration preservation**: Your .env settings are maintained
- **Rollback capability**: Keeps backup of previous version
- **Git integration**: Works with git repositories for version control
- **Standalone mode**: Can update even without git

### Update Process

1. Checks for new versions
2. Backs up current configuration (.env → .env.backup)
3. Saves any local changes
4. Downloads and applies updates
5. Rebuilds the executable
6. Preserves your configuration
7. Shows changelog of new features

## Troubleshooting

### "tools" command not found
- Restart your terminal after installation
- Verify PATH was updated: `echo %PATH%`
- Check if tools.exe exists in the installation directory

### Permission denied
- Make sure to run installation scripts as Administrator
- Check file permissions on tools.exe

### Configuration not loading
- Check if .env file exists in the tools directory
- Verify environment variables: `tools config env`
- Check if variables are set correctly: `echo %PGHOST%`

### Build errors
- Ensure Go is installed: `go version`
- Update dependencies: `go mod download`
- Clean and rebuild: `go clean && go build -o tools.exe .`

## Support

For issues or questions, check:
- `tools --help` for command documentation
- `tools config all` to verify configuration
- README.md for usage examples