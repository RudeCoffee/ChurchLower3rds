@echo off
setlocal
cd /d "%~dp0"

echo ===================================================
echo Building Church Lower Thirds for Windows (64-bit)
echo ===================================================

set GOOS=windows
set GOARCH=amd64

:: Check for GCC
where gcc >nul 2>nul
if %ERRORLEVEL% NEQ 0 goto NO_GCC

:: Check for Vosk CGO headers / library in root folder
if exist "%CD%\libvosk.dll" (
    set CGO_CFLAGS=-I"%CD%"
    set CGO_LDFLAGS=-L"%CD%" -lvosk
    goto WITH_VOSK
)

:: Check for Vosk CGO headers / library in ./vosk/lib/ folder
if exist "%CD%\vosk\lib\vosk_api.h" (
    set CGO_CFLAGS=-I"%CD%\vosk\lib"
    set CGO_LDFLAGS=-L"%CD%\vosk\lib" -lvosk
    goto WITH_VOSK
)

:NO_VOSK
echo [NOTE] libvosk.dll or vosk_api.h not found in project root or ./vosk/lib/.
echo Compiling in Stub Mode (Voice Assistant disabled).
echo To enable Voice Assistant:
echo   1. Download libvosk.dll and vosk_api.h from https://github.com/alphacep/vosk-api/releases
echo   2. Place libvosk.dll and vosk_api.h in the project folder.
echo   3. Re-run build_windows.bat.
echo.
goto STUB_BUILD

:NO_GCC
echo [NOTE] GCC compiler not found in PATH.
echo Compiling in Stub Mode (Voice Assistant disabled).
echo To enable Voice Assistant, install TDM-GCC or MinGW.
echo.
goto STUB_BUILD

:WITH_VOSK
echo GCC compiler and Vosk library detected. Compiling with full CGO Voice Assistant...
set CGO_ENABLED=1
go build -o church-lower-thirds.exe
if %ERRORLEVEL% EQU 0 goto CHECK_RESULT
echo.
echo [WARNING] CGO build failed. Falling back to Stub Mode...

:STUB_BUILD
set CGO_ENABLED=0
go build -tags stub -o church-lower-thirds.exe

:CHECK_RESULT
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Build failed! Please ensure Go is properly installed.
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo ===================================================
echo [SUCCESS] church-lower-thirds.exe compiled successfully!
echo ===================================================
pause
