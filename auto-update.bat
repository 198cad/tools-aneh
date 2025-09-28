@echo off
:: Auto-update checker for Development Tools
:: This can be added to Windows Task Scheduler to run periodically

cd /d "C:\Users\mafiaboy\Projects\tools"

:: Check for updates silently
tools.exe update --check > update_check.log 2>&1

:: Check if updates are available
findstr /C:"Updates available" update_check.log >nul
if %errorLevel% == 0 (
    echo Updates available! Run 'tools update' to install.

    :: Optional: Show notification (Windows 10/11)
    powershell -Command "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('Updates available for Development Tools!`nRun: tools update', 'Tools Update', 'OK', 'Information')" >nul 2>&1
)

:: Clean up
del update_check.log >nul 2>&1
