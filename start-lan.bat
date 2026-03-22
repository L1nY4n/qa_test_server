@echo off
setlocal
set "SCRIPT_DIR=%~dp0"

powershell -NoProfile -ExecutionPolicy Bypass -File "%SCRIPT_DIR%start-lan.ps1" %*
if errorlevel 1 (
  echo.
  echo [ERROR] LAN startup failed.
  exit /b %errorlevel%
)

echo.
echo [OK] LAN startup complete.
