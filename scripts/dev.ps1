# BLOCK_SCRIPT_DEV_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [dev] Starting Campus OS backend database service..." -ForegroundColor Green
docker compose up -d campus-db

Write-Host "==> [dev] Starting TypeScript DB Layer in development mode..." -ForegroundColor Green
Push-Location "$ScriptDir\db-layer"; pnpm run dev; Pop-Location
