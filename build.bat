@echo off
setlocal

pushd "%~dp0" >nul || (
    echo Error: cannot enter the project directory.
    exit /b 1
)

set "GO_CMD="
where go.exe >nul 2>&1 && set "GO_CMD=go.exe"

if not defined GO_CMD (
    if exist "C:\Program Files\Go\bin\go.exe" (
        set "GO_CMD=C:\Program Files\Go\bin\go.exe"
    )
)

if not defined GO_CMD (
    echo Error: Go was not found in PATH or C:\Program Files\Go\bin.
    popd >nul
    exit /b 1
)

echo Building jsondiff.exe...
"%GO_CMD%" build -o jsondiff.exe .
set "BUILD_EXIT=%ERRORLEVEL%"

if not "%BUILD_EXIT%"=="0" (
    echo Build failed with exit code %BUILD_EXIT%.
    popd >nul
    exit /b %BUILD_EXIT%
)

echo Build complete: %CD%\jsondiff.exe
popd >nul
exit /b 0
