@echo off
setlocal enabledelayedexpansion

echo ===============================================
echo  Development Tools - Update Script
echo ===============================================
echo.

:: Get current directory
set INSTALL_DIR=%~dp0
set INSTALL_DIR=%INSTALL_DIR:~0,-1%

cd /d "%INSTALL_DIR%"

:: Check if git is available
where git >nul 2>&1
if %errorLevel% neq 0 (
    echo Git is not installed or not in PATH!
    echo Please install Git from https://git-scm.com/
    goto :BUILD_ONLY
)

:: Check if it's a git repository
if not exist ".git" (
    echo This is not a git repository.
    echo Performing direct build update...
    goto :BUILD_ONLY
)

:: Fetch latest changes
echo Checking for updates...
git fetch origin >nul 2>&1

:: Check if there are updates
for /f "delims=" %%i in ('git rev-list HEAD...origin/main --count 2^>nul') do set BEHIND=%%i
if "%BEHIND%"=="" set BEHIND=0

if %BEHIND% equ 0 (
    echo You are already running the latest version!
    set /p FORCE="Do you want to force rebuild? (y/n): "
    if /i "!FORCE!" neq "y" goto :END
)

:: Show what's new
if %BEHIND% gtr 0 (
    echo.
    echo Found %BEHIND% new updates:
    git log --oneline HEAD..origin/main --max-count=10
    echo.
)

:: Backup current configuration
echo Backing up configuration...
if exist ".env" (
    copy /Y ".env" ".env.backup" >nul
    echo Configuration backed up to .env.backup
)

:: Stash any local changes
echo Saving local changes...
git stash save "Auto-stash before update %date% %time%" >nul 2>&1

:: Pull latest changes
echo Downloading updates...
git pull origin main
if %errorLevel% neq 0 (
    echo Error pulling updates. Trying to resolve...
    git reset --hard origin/main
    if %errorLevel% neq 0 (
        echo Failed to update from repository.
        goto :BUILD_ONLY
    )
)

:: Restore configuration
if exist ".env.backup" (
    if not exist ".env" (
        copy /Y ".env.backup" ".env" >nul
        echo Configuration restored.
    )
)

:BUILD_ONLY
:: Install/Update dependencies
echo Updating dependencies...
go mod download
if %errorLevel% neq 0 (
    echo Warning: Could not update dependencies
)

:: Build new version
echo Building new version...
go build -o tools.exe.new .
if %errorLevel% neq 0 (
    echo Build failed!
    pause
    exit /b 1
)

:: Replace the executable
echo Installing new version...

:: Create replacement script
echo @echo off > replace_tools.bat
echo :RETRY >> replace_tools.bat
echo timeout /t 1 /nobreak ^> nul >> replace_tools.bat
echo move /Y tools.exe tools.exe.old ^>nul 2^>^&1 >> replace_tools.bat
echo if errorlevel 1 goto :RETRY >> replace_tools.bat
echo move /Y tools.exe.new tools.exe >> replace_tools.bat
echo if exist tools.exe.old del tools.exe.old >> replace_tools.bat
echo echo Update completed successfully! >> replace_tools.bat
echo timeout /t 2 /nobreak ^> nul >> replace_tools.bat
echo del "%%~f0" >> replace_tools.bat

:: Run replacement in background
start /b cmd /c replace_tools.bat

echo.
echo ===============================================
echo  Update Complete!
echo ===============================================
echo.
echo The tools have been updated to the latest version.
echo Please restart your terminal to use the new version.
echo.

:: Show version
tools.exe.new --version 2>nul
echo.

:END
pause
