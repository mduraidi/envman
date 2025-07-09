@echo off
REM =============================================================
REM Build script for envman (Portable SDK Environment Manager)
REM -------------------------------------------------------------
REM Builds envman.exe using a portable Go toolchain. Also copies
REM documentation and templates, then packages the result as a ZIP.
REM -------------------------------------------------------------
REM Prerequisites:
REM   - Go 1.22+ extracted at %GOROOT%
REM   - Script must run from the envman source folder
REM =============================================================

REM Configure environment
setlocal
set GO111MODULE=on
set GOROOT=D:\envman\sdks\go\1.24.5
set PATH=%GOROOT%\bin;%PATH%

REM Validate Go binary presence
if not exist "%GOROOT%\bin\go.exe" (
    echo [ERROR] go.exe not found in %GOROOT%\bin
    exit /b 1
)

REM Show Go version
echo ----------------------------------------
echo Using Go version:
go version
echo ----------------------------------------

REM Download modules and vendorize
echo Tidying and vendoring modules...
go mod tidy
go mod vendor

REM Prepare dist folder
if not exist dist mkdir dist

REM Clean previous binary
del /f /q envman.exe 2>nul

REM Build executable
echo Building envman.exe ...
go build -mod=vendor -o dist\envman.exe
if errorlevel 1 (
    echo [ERROR] Build failed. Check your Go code for issues.
    exit /b 1
)
echo Build succeeded.

REM Copy documentation and templates
echo Copying documentation and assets...
if exist README.md copy /y README.md dist\README.md >nul
if exist templates (
    mkdir dist\templates >nul 2>nul
    xcopy /e /i /y /q templates dist\templates >nul
)

REM (Optional) Copy example SDK layout
REM if exist envman_sdks (
REM     mkdir dist\envman_sdks >nul 2>nul
REM     xcopy /e /i /y /q envman_sdks dist\envman_sdks >nul
REM )

REM Final packaging into a ZIP file
echo Packaging dist...
del /f /q envman_bundle.zip 2>nul
powershell -Command "Compress-Archive -Path dist -DestinationPath envman_bundle.zip"

echo.
echo ✅ Done. Your envman bundle is ready as envman_bundle.zip
echo Contents:
echo   - envman.exe
echo   - README.md (if present)
echo   - templates\ (if present)
echo.
endlocal