# Development Tools - PowerShell Update Script

Write-Host "===============================================" -ForegroundColor Cyan
Write-Host " Development Tools - Update Script" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""

# Get installation directory
$InstallDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $InstallDir

# Function to check if process is running
function Test-ProcessRunning {
    param([string]$ProcessName)
    $process = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue
    return $process -ne $null
}

# Check if git is available
$gitAvailable = Get-Command git -ErrorAction SilentlyContinue
if (-not $gitAvailable) {
    Write-Host "Git is not installed or not in PATH!" -ForegroundColor Red
    Write-Host "Please install Git from https://git-scm.com/" -ForegroundColor Yellow
    Write-Host "Performing direct build update..." -ForegroundColor Yellow
    $gitUpdate = $false
} else {
    $gitUpdate = Test-Path ".git"
}

if ($gitUpdate) {
    # Fetch latest changes
    Write-Host "Checking for updates..." -ForegroundColor Yellow
    $fetchResult = git fetch origin 2>&1

    # Check if there are updates
    $behind = git rev-list HEAD...origin/main --count 2>$null
    if (-not $behind) { $behind = 0 }

    if ($behind -eq 0) {
        Write-Host "You are already running the latest version!" -ForegroundColor Green
        $force = Read-Host "Do you want to force rebuild? (y/n)"
        if ($force -ne 'y') {
            Write-Host "No updates needed." -ForegroundColor Green
            Read-Host "Press Enter to exit"
            exit 0
        }
    } else {
        # Show what's new
        Write-Host ""
        Write-Host "Found $behind new updates:" -ForegroundColor Green
        git log --oneline HEAD..origin/main --max-count=10
        Write-Host ""
    }

    # Backup current configuration
    Write-Host "Backing up configuration..." -ForegroundColor Yellow
    if (Test-Path ".env") {
        Copy-Item ".env" ".env.backup" -Force
        Write-Host "Configuration backed up to .env.backup" -ForegroundColor Green
    }

    # Stash any local changes
    Write-Host "Saving local changes..." -ForegroundColor Yellow
    $stashResult = git stash save "Auto-stash before update $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" 2>&1
    $hasStash = $stashResult -notlike "*No local changes*"

    # Pull latest changes
    Write-Host "Downloading updates..." -ForegroundColor Yellow
    $pullResult = git pull origin main 2>&1

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error pulling updates. Trying to resolve..." -ForegroundColor Yellow
        git reset --hard origin/main
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Failed to update from repository." -ForegroundColor Red
            $gitUpdate = $false
        }
    } else {
        Write-Host $pullResult -ForegroundColor Green
    }

    # Restore configuration
    if (Test-Path ".env.backup") {
        if (-not (Test-Path ".env")) {
            Copy-Item ".env.backup" ".env" -Force
            Write-Host "Configuration restored." -ForegroundColor Green
        }
    }

    # Restore stashed changes if any
    if ($hasStash) {
        Write-Host "Restoring local changes..." -ForegroundColor Yellow
        git stash pop 2>&1 | Out-Null
    }
}

# Install/Update dependencies
Write-Host "Updating dependencies..." -ForegroundColor Yellow
$modResult = go mod download 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "Warning: Could not update dependencies" -ForegroundColor Yellow
}

# Build new version
Write-Host "Building new version..." -ForegroundColor Yellow
$buildResult = go build -o tools.exe.new . 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    Write-Host $buildResult
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "Build successful!" -ForegroundColor Green

# Replace the executable
Write-Host "Installing new version..." -ForegroundColor Yellow

# Check if tools.exe is running
if (Test-ProcessRunning "tools") {
    Write-Host "tools.exe is currently running. Scheduling update..." -ForegroundColor Yellow

    # Create a scheduled task to replace the file
    $action = New-ScheduledTaskAction -Execute "powershell.exe" -Argument "-WindowStyle Hidden -Command `"Stop-Process -Name tools -Force -ErrorAction SilentlyContinue; Start-Sleep -Seconds 2; Move-Item -Path '$InstallDir\tools.exe' -Destination '$InstallDir\tools.exe.old' -Force; Move-Item -Path '$InstallDir\tools.exe.new' -Destination '$InstallDir\tools.exe' -Force; Remove-Item -Path '$InstallDir\tools.exe.old' -Force -ErrorAction SilentlyContinue`""
    $trigger = New-ScheduledTaskTrigger -Once -At (Get-Date).AddSeconds(5)
    $principal = New-ScheduledTaskPrincipal -UserId $env:USERNAME -LogonType Interactive -RunLevel Highest
    $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

    Register-ScheduledTask -TaskName "UpdateDevTools" -Action $action -Trigger $trigger -Principal $principal -Settings $settings -Force | Out-Null

    Write-Host "Update scheduled. The new version will be installed shortly." -ForegroundColor Green
    Write-Host "Please close any running instances of tools.exe" -ForegroundColor Yellow
} else {
    # Direct replacement
    try {
        if (Test-Path "tools.exe") {
            Move-Item -Path "tools.exe" -Destination "tools.exe.old" -Force
        }
        Move-Item -Path "tools.exe.new" -Destination "tools.exe" -Force

        if (Test-Path "tools.exe.old") {
            Remove-Item -Path "tools.exe.old" -Force
        }

        Write-Host "Update installed successfully!" -ForegroundColor Green
    } catch {
        Write-Host "Error replacing executable: $_" -ForegroundColor Red

        # Try to restore old version
        if (Test-Path "tools.exe.old") {
            Move-Item -Path "tools.exe.old" -Destination "tools.exe" -Force
        }

        Read-Host "Press Enter to exit"
        exit 1
    }
}

Write-Host ""
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host " Update Complete!" -ForegroundColor Green
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""

# Show new version
try {
    $version = & "$InstallDir\tools.exe" --version 2>&1
    Write-Host "New version: $version" -ForegroundColor Green
} catch {
    Write-Host "Version information not available" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "The tools have been updated to the latest version." -ForegroundColor Green
Write-Host "Please restart your terminal to use the new version." -ForegroundColor Yellow
Write-Host ""

# Clean up scheduled task if it exists
Unregister-ScheduledTask -TaskName "UpdateDevTools" -Confirm:$false -ErrorAction SilentlyContinue

Read-Host "Press Enter to exit"
