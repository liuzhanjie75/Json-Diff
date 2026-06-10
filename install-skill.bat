@echo off
setlocal EnableExtensions

pushd "%~dp0" >nul || (
    echo Error: cannot enter the project directory.
    exit /b 1
)

set "GO_CMD="
where go.exe >nul 2>&1 && set "GO_CMD=go.exe"
if not defined GO_CMD if exist "C:\Program Files\Go\bin\go.exe" set "GO_CMD=C:\Program Files\Go\bin\go.exe"

if not defined GO_CMD (
    echo Error: Go was not found in PATH or C:\Program Files\Go\bin.
    popd >nul
    exit /b 1
)

if defined CODEX_HOME (
    set "CODEX_ROOT=%CODEX_HOME%"
) else (
    set "CODEX_ROOT=%USERPROFILE%\.codex"
)

set "SKILLS_DIR=%CODEX_ROOT%\skills"
set "TARGET_DIR=%SKILLS_DIR%\jsondiff"
set "STAGING_DIR=%SKILLS_DIR%\.jsondiff-install-%RANDOM%-%RANDOM%"
set "BACKUP_DIR=%SKILLS_DIR%\.jsondiff-backup-%RANDOM%-%RANDOM%"
set "GOCACHE=%CD%\.gocache-skill-install"

if not exist "%SKILLS_DIR%" mkdir "%SKILLS_DIR%" || goto :fail
mkdir "%STAGING_DIR%" || goto :fail

xcopy "skill\jsondiff\*" "%STAGING_DIR%\" /E /I /Q /Y >nul || goto :fail
if not exist "%STAGING_DIR%\bin" mkdir "%STAGING_DIR%\bin" || goto :fail

echo Building jsondiff for the skill...
"%GO_CMD%" build -o "%STAGING_DIR%\bin\jsondiff.exe" .
if errorlevel 1 goto :fail

if exist "%TARGET_DIR%" (
    move "%TARGET_DIR%" "%BACKUP_DIR%" >nul || goto :fail
)

move "%STAGING_DIR%" "%TARGET_DIR%" >nul || goto :rollback

if exist "%BACKUP_DIR%" rmdir /S /Q "%BACKUP_DIR%"
if exist "%GOCACHE%" rmdir /S /Q "%GOCACHE%"
echo Installed JSON Diff skill: %TARGET_DIR%
echo Restart Codex or start a new session to discover the skill.
popd >nul
exit /b 0

:rollback
if exist "%BACKUP_DIR%" move "%BACKUP_DIR%" "%TARGET_DIR%" >nul

:fail
set "INSTALL_EXIT=%ERRORLEVEL%"
if "%INSTALL_EXIT%"=="0" set "INSTALL_EXIT=1"
if defined STAGING_DIR if exist "%STAGING_DIR%" rmdir /S /Q "%STAGING_DIR%"
if defined GOCACHE if exist "%GOCACHE%" rmdir /S /Q "%GOCACHE%"
echo Skill installation failed.
popd >nul
exit /b %INSTALL_EXIT%
