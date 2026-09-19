@echo off
setlocal
cd /d "%~dp0"

echo ===================================================
echo Building Church Lower Thirds for Windows (64-bit)
echo ===================================================

set GOOS=windows
set GOARCH=amd64

where gcc >nul 2>nul
if %ERRORLEVEL% EQU 0 goto WITH_GCC
goto WITHOUT_GCC

:WITH_GCC
echo GCC compiler detected. Compiling with CGO support for speech recognition...
set CGO_ENABLED=1
go build -o church-lower-thirds.exe
goto CHECK_RESULT

:WITHOUT_GCC
echo [WARNING] GCC compiler not detected in PATH!
echo Compiling in Stub Mode (Voice recognition disabled).
echo To enable Voice Assistant, install TDM-GCC or MinGW and ensure libvosk.dll is present.
echo.
set CGO_ENABLED=0
go build -tags stub -o church-lower-thirds.exe
goto CHECK_RESULT

:CHECK_RESULT
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Build failed! Please ensure Go is properly installed.
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo [SUCCESS] church-lower-thirds.exe compiled successfully!
pause
