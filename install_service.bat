@echo off
setlocal
:: Ensure administrative privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] This script must be run as Administrator!
    echo Right-click install_service.bat and select "Run as administrator".
    pause
    exit /b 1
)

echo ===================================================
echo Installing Church Lower Thirds Windows Service
echo ===================================================

set SERVICE_NAME=ChurchLowerThirds
set DISPLAY_NAME=Church Lower Thirds Service
set BIN_PATH=%~dp0church-lower-thirds.exe

if not exist "%BIN_PATH%" (
    echo [ERROR] %BIN_PATH% not found!
    echo Please run build_windows.bat first.
    pause
    exit /b 1
)

echo Registering Windows Service '%SERVICE_NAME%'...
sc create %SERVICE_NAME% binPath= "\"%BIN_PATH%\" -service" start= auto displayName= "%DISPLAY_NAME%"

if %ERRORLEVEL% equ 0 (
    echo Setting service description...
    sc description %SERVICE_NAME% "Church Lower Thirds real-time Bible display and voice assistant server."
    echo Starting service...
    sc start %SERVICE_NAME%
    echo ===================================================
    echo [SUCCESS] Church Lower Thirds Service installed and started!
    echo It will now start automatically whenever Windows boots.
    echo Control Panel: http://localhost:8080/client.html
    echo OBS Browser Source: http://localhost:8080/obs.html
    echo ===================================================
) else (
    echo [ERROR] Failed to create service. Error code: %ERRORLEVEL%
)

pause
