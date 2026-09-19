@echo off
setlocal
echo ===================================================
echo Building Church Lower Thirds for Windows Service
echo ===================================================

set CGO_ENABLED=1

echo Compiling executable...
go build -o church-lower-thirds.exe .

if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Build failed. Please ensure Go and TDM-GCC/MinGW are installed.
    exit /b %ERRORLEVEL%
)

echo [SUCCESS] church-lower-thirds.exe compiled successfully!
pause
