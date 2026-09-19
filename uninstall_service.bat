@echo off
setlocal
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] This script must be run as Administrator!
    echo Right-click uninstall_service.bat and select "Run as administrator".
    pause
    exit /b 1
)

echo ===================================================
echo Uninstalling Church Lower Thirds Windows Service
echo ===================================================

set SERVICE_NAME=ChurchLowerThirds

echo Stopping service '%SERVICE_NAME%'...
sc stop %SERVICE_NAME%

echo Deleting service '%SERVICE_NAME%'...
sc delete %SERVICE_NAME%

if %ERRORLEVEL% equ 0 (
    echo [SUCCESS] Church Lower Thirds Service uninstalled successfully!
) else (
    echo [ERROR] Failed to delete service. Error code: %ERRORLEVEL%
)

pause
