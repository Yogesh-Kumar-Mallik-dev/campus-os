# BLOCK_SCRIPT_DEV_DB_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [dev:db] Launching TypeScript DB Layer with tsx watcher..." -ForegroundColor Green
Push-Location "$ScriptDir\db-layer"
try {
    pnpm run dev
} finally {
    Pop-Location
}
