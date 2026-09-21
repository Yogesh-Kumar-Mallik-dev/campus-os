# BLOCK_SCRIPT_DEV_CLIENT_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [dev:client] Launching Go Fyne Native Client..." -ForegroundColor Green
Push-Location "$ScriptDir\apps\client"
try {
    go run . $args
} finally {
    Pop-Location
}
