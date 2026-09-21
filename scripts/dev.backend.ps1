# BLOCK_SCRIPT_DEV_BACKEND_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [dev:backend] Launching Go Logical Backend on :8080..." -ForegroundColor Green
Push-Location "$ScriptDir\backend"
try {
    go run ./cmd/server $args
} finally {
    Pop-Location
}
