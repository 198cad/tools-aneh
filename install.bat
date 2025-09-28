@echo off
setlocal enabledelayedexpansion

echo ===============================================
echo  Development Tools CLI - Installation Script
echo ===============================================
echo.

:: Check if running as administrator
net session >nul 2>&1
if %errorLevel% == 0 (
    echo Running with Administrator privileges...
) else (
    echo This script requires Administrator privileges!
    echo Please run as Administrator.
    pause
    exit /b 1
)

:: Get current directory
set INSTALL_DIR=%~dp0
set INSTALL_DIR=%INSTALL_DIR:~0,-1%

echo Installation directory: %INSTALL_DIR%
echo.

:: Build the executable
echo Building tools.exe...
go build -o tools.exe .
if %errorLevel% neq 0 (
    echo Failed to build tools.exe
    pause
    exit /b 1
)
echo Build successful!
echo.

:: Check if directory is already in PATH
echo Checking PATH environment variable...
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if %errorLevel% == 0 (
    echo Directory already in PATH!
) else (
    echo Adding to system PATH...

    :: Add to system PATH
    for /f "tokens=2*" %%a in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v Path 2^>nul') do set "OLDPATH=%%b"

    setx /M PATH "!OLDPATH!;%INSTALL_DIR%" >nul 2>&1

    if !errorLevel! == 0 (
        echo Successfully added to system PATH!
    ) else (
        echo Failed to add to system PATH. Trying user PATH...

        :: Try adding to user PATH instead
        for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USERPATH=%%b"
        setx PATH "!USERPATH!;%INSTALL_DIR%" >nul 2>&1

        if !errorLevel! == 0 (
            echo Successfully added to user PATH!
        ) else (
            echo Failed to update PATH variable
            pause
            exit /b 1
        )
    )
)

:: Copy .env.example if .env doesn't exist
if not exist "%INSTALL_DIR%\.env" (
    if exist "%INSTALL_DIR%\.env.example" (
        echo Creating .env from .env.example...
        copy "%INSTALL_DIR%\.env.example" "%INSTALL_DIR%\.env" >nul
        echo Please edit .env file with your configuration
    )
)

echo.
echo ===============================================
echo  Installation Complete!
echo ===============================================
echo.
echo The 'tools' command is now available globally.
echo.
echo IMPORTANT: You need to restart your terminal or run:
echo   refreshenv
echo.
echo Usage examples:
echo   tools --help
echo   tools config all
echo   tools db list
echo   tools rabbit queues
echo   tools minio buckets
echo.
echo Configuration:
echo   Edit the .env file in: %INSTALL_DIR%
echo.

pause
