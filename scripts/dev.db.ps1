# BLOCK_SCRIPT_DEV_DB_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
if (-not $env:SOCKET_PATH) { $env:SOCKET_PATH = "$env:TEMP\campus-os-dev.sock" }

Write-Host "==> [dev:db] Launching TypeScript DB Layer with tsx watcher on $env:SOCKET_PATH..." -ForegroundColor Green
Push-Location "$ScriptDir\db-layer"
try {
    pnpm run dev @args
} finally {
    Pop-Location
}
