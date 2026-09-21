# BLOCK_SCRIPT_DEV_CLIENT_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
if (-not $env:CAMPUS_BACKEND_URL) { $env:CAMPUS_BACKEND_URL = "http://localhost:8080" }

Write-Host "==> [dev:client] Launching Go Fyne Native Client (Backend: $env:CAMPUS_BACKEND_URL)..." -ForegroundColor Green
Push-Location "$ScriptDir\apps\client"
try {
    go run . $args
} finally {
    Pop-Location
}
