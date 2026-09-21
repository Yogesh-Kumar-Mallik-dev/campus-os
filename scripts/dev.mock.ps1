# BLOCK_SCRIPT_DEV_MOCK_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$SocketPath = "$env:TEMP\campus-os-mock.sock"

Write-Host "========================================================================" -ForegroundColor Cyan
Write-Host "  CAMPUS OS: ZERO-DEPENDENCY MOCK DEVELOPMENT ENVIRONMENT               " -ForegroundColor Cyan
Write-Host "========================================================================" -ForegroundColor Cyan
Write-Host "  Mock Socket Path: unix://$SocketPath" -ForegroundColor Yellow
Write-Host "  Backend API:      http://localhost:8080" -ForegroundColor Yellow
Write-Host "========================================================================" -ForegroundColor Cyan

if (Test-Path $SocketPath) {
    Remove-Item $SocketPath -Force
}

$dbProcess = $null
$backendProcess = $null

try {
    Write-Host "==> [dev:mock] Starting in-memory mock persistence server..." -ForegroundColor Green
    $env:SOCKET_PATH = $SocketPath
    $dbProcess = Start-Process pnpm -ArgumentList "--filter @campus-os/db-layer exec tsx src/server.mock.ts" -WorkingDirectory "$ScriptDir\db-layer" -PassThru

    Write-Host "==> [dev:mock] Waiting for mock socket to initialize..." -ForegroundColor Green
    $ready = $false
    for ($i = 0; $i -lt 30; $i++) {
        if (Test-Path $SocketPath) {
            $ready = $true
            break
        }
        Start-Sleep -Milliseconds 200
    }

    if (-not $ready) {
        Write-Error "Mock persistence server socket failed to bind at $SocketPath"
        exit 1
    }

    Write-Host "==> [dev:mock] Starting Go logical backend..." -ForegroundColor Green
    $env:DB_SOCKET_PATH = $SocketPath
    $env:PORT = "8080"
    $backendProcess = Start-Process go -ArgumentList "run ./cmd/server" -WorkingDirectory "$ScriptDir\backend" -PassThru

    Write-Host "==> [dev:mock] Waiting for HTTP API on http://localhost:8080/healthz..." -ForegroundColor Green
    for ($i = 0; $i -lt 30; $i++) {
        try {
            $resp = Invoke-WebRequest -Uri "http://localhost:8080/healthz" -UseBasicParsing -TimeoutSec 1
            if ($resp.StatusCode -eq 200) { break }
        } catch {
            Start-Sleep -Milliseconds 200
        }
    }

    Write-Host "==> [dev:mock] Launching Native Fyne Client..." -ForegroundColor Green
    Push-Location "$ScriptDir\apps\client"
    go run .
    Pop-Location
} finally {
    Write-Host "`n==> [dev:mock] Cleaning up mock development processes..." -ForegroundColor Yellow
    if ($backendProcess -and -not $backendProcess.HasExited) { Stop-Process -Id $backendProcess.Id -Force }
    if ($dbProcess -and -not $dbProcess.HasExited) { Stop-Process -Id $dbProcess.Id -Force }
    if (Test-Path $SocketPath) { Remove-Item $SocketPath -Force }
    Write-Host "==> [dev:mock] Clean shutdown complete." -ForegroundColor Green
}
