# Development Tools CLI - PowerShell Installation Script

Write-Host "===============================================" -ForegroundColor Cyan
Write-Host " Development Tools CLI - Installation Script" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")

if (-not $isAdmin) {
    Write-Host "This script requires Administrator privileges!" -ForegroundColor Red
    Write-Host "Restarting as Administrator..." -ForegroundColor Yellow

    # Restart script as Administrator
    Start-Process PowerShell -Verb RunAs -ArgumentList "-ExecutionPolicy Bypass -File `"$PSCommandPath`""
    exit
}

# Get installation directory
$InstallDir = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "Installation directory: $InstallDir" -ForegroundColor Green
Write-Host ""

# Build the executable
Write-Host "Building tools.exe..." -ForegroundColor Yellow
Set-Location $InstallDir
$buildResult = & go build -o tools.exe . 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build tools.exe" -ForegroundColor Red
    Write-Host $buildResult
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "Build successful!" -ForegroundColor Green
Write-Host ""

# Function to add to PATH
function Add-ToPath {
    param([string]$Path, [string]$Scope)

    $currentPath = [Environment]::GetEnvironmentVariable("Path", $Scope)

    if ($currentPath -notlike "*$Path*") {
        $newPath = $currentPath + ";" + $Path
        [Environment]::SetEnvironmentVariable("Path", $newPath, $Scope)
        return $true
    }
    return $false
}

# Try to add to system PATH
Write-Host "Updating PATH environment variable..." -ForegroundColor Yellow

try {
    if (Add-ToPath -Path $InstallDir -Scope "Machine") {
        Write-Host "Successfully added to system PATH!" -ForegroundColor Green
    } else {
        Write-Host "Directory already in system PATH" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Failed to add to system PATH, trying user PATH..." -ForegroundColor Yellow

    try {
        if (Add-ToPath -Path $InstallDir -Scope "User") {
            Write-Host "Successfully added to user PATH!" -ForegroundColor Green
        } else {
            Write-Host "Directory already in user PATH" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "Failed to update PATH variable" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        Read-Host "Press Enter to exit"
        exit 1
    }
}

# Refresh PATH for current session
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Copy .env.example if .env doesn't exist
$envFile = Join-Path $InstallDir ".env"
$envExampleFile = Join-Path $InstallDir ".env.example"

if (-not (Test-Path $envFile)) {
    if (Test-Path $envExampleFile) {
        Write-Host "Creating .env from .env.example..." -ForegroundColor Yellow
        Copy-Item $envExampleFile $envFile
        Write-Host "Please edit .env file with your configuration" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host " Installation Complete!" -ForegroundColor Green
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "The 'tools' command is now available globally." -ForegroundColor Green
Write-Host ""
Write-Host "IMPORTANT: You need to restart your terminal/PowerShell" -ForegroundColor Yellow
Write-Host ""
Write-Host "Usage examples:" -ForegroundColor Cyan
Write-Host "  tools --help"
Write-Host "  tools config all"
Write-Host "  tools db list"
Write-Host "  tools rabbit queues"
Write-Host "  tools minio buckets"
Write-Host ""
Write-Host "Configuration:" -ForegroundColor Cyan
Write-Host "  Edit the .env file in: $InstallDir" -ForegroundColor Yellow
Write-Host ""

# Test if tools command works
Write-Host "Testing installation..." -ForegroundColor Yellow
try {
    $testResult = & "$InstallDir\tools.exe" --version 2>&1
    Write-Host "✓ Tools version: $testResult" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to run tools.exe" -ForegroundColor Red
}

Write-Host ""
Read-Host "Press Enter to exit"
