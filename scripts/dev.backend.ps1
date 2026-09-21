# BLOCK_SCRIPT_DEV_BACKEND_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
if (-not $env:DB_SOCKET_PATH) { $env:DB_SOCKET_PATH = "$env:TEMP\campus-os-dev.sock" }
if (-not $env:PORT) { $env:PORT = "8080" }
if (-not $env:ENV) { $env:ENV = "development" }

Write-Host "==> [dev:backend] Launching Go Logical Backend on :$env:PORT (Socket: $env:DB_SOCKET_PATH)..." -ForegroundColor Green
Push-Location "$ScriptDir\backend"
try {
    go run ./cmd/server $args
} finally {
    Pop-Location
}
