# BLOCK_SCRIPT_DEPS_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [deps] Installing pnpm workspace dependencies..." -ForegroundColor Green
pnpm install

Write-Host "==> [deps] Downloading Go module dependencies..." -ForegroundColor Green
Push-Location "$ScriptDir\backend"; go mod download; Pop-Location
Push-Location "$ScriptDir\apps\client"; go mod download; Pop-Location

Write-Host "==> [deps] All dependencies installed." -ForegroundColor Green
