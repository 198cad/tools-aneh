@echo off
setlocal enabledelayedexpansion

echo ===============================================
echo  Development Tools CLI - Uninstall Script
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

echo Uninstalling from: %INSTALL_DIR%
echo.

:: Remove from system PATH
echo Removing from PATH environment variable...

:: Get current system PATH
for /f "tokens=2*" %%a in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v Path 2^>nul') do set "SYSPATH=%%b"

:: Remove install directory from PATH
set "NEWPATH=!SYSPATH!"
set "NEWPATH=!NEWPATH:;%INSTALL_DIR%=!"
set "NEWPATH=!NEWPATH:%INSTALL_DIR%;=!"

:: Update system PATH
reg add "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v Path /t REG_EXPAND_SZ /d "!NEWPATH!" /f >nul 2>&1

if !errorLevel! == 0 (
    echo Successfully removed from system PATH!
) else (
    echo Could not remove from system PATH, checking user PATH...

    :: Try user PATH
    for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USERPATH=%%b"

    :: Remove install directory from user PATH
    set "NEWUSERPATH=!USERPATH!"
    set "NEWUSERPATH=!NEWUSERPATH:;%INSTALL_DIR%=!"
    set "NEWUSERPATH=!NEWUSERPATH:%INSTALL_DIR%;=!"

    setx PATH "!NEWUSERPATH!" >nul 2>&1

    if !errorLevel! == 0 (
        echo Successfully removed from user PATH!
    )
)

echo.
echo ===============================================
echo  Uninstallation Complete!
echo ===============================================
echo.
echo The 'tools' command has been removed from PATH.
echo You need to restart your terminal for changes to take effect.
echo.
echo The files are still in: %INSTALL_DIR%
echo You can manually delete this folder if desired.
echo.

pause
