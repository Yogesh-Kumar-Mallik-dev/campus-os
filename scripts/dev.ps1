# BLOCK_SCRIPT_DEV_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$Subcommand = if ($args.Count -gt 0) { $args[0] } else { "" }
$SubArgs = if ($args.Count -gt 1) { $args[1..($args.Count - 1)] } else { @() }

switch ($Subcommand) {
    "mock" {
        & "$ScriptDir\scripts\dev.mock.ps1" @SubArgs
    }
    "backend" {
        & "$ScriptDir\scripts\dev.backend.ps1" @SubArgs
    }
    "db" {
        & "$ScriptDir\scripts\dev.db.ps1" @SubArgs
    }
    "client" {
        & "$ScriptDir\scripts\dev.client.ps1" @SubArgs
    }
    default {
        $SocketPath = if ($env:SOCKET_PATH) { $env:SOCKET_PATH } else { "$env:TEMP\campus-os-dev.sock" }
        $Port = if ($env:PORT) { $env:PORT } else { "8080" }
        $BackendUrl = if ($env:CAMPUS_BACKEND_URL) { $env:CAMPUS_BACKEND_URL } else { "http://localhost:$Port" }

        Write-Host "========================================================================" -ForegroundColor Cyan
        Write-Host "  CAMPUS OS: LOCAL FULL-STACK DEVELOPMENT ENVIRONMENT                  " -ForegroundColor Cyan
        Write-Host "========================================================================" -ForegroundColor Cyan
        Write-Host "  Persistence Socket: unix://$SocketPath" -ForegroundColor Yellow
        Write-Host "  Backend API:        $BackendUrl" -ForegroundColor Yellow
        Write-Host "  Client UI:          Native Fyne (Go)" -ForegroundColor Yellow
        Write-Host "========================================================================" -ForegroundColor Cyan

        if (Test-Path $SocketPath) {
            Remove-Item $SocketPath -Force
        }

        $dbProcess = $null
        $backendProcess = $null
        $clientProcess = $null

        try {
            Write-Host "==> [dev] Ensuring Campus OS PostgreSQL database is running..." -ForegroundColor Green
            docker compose up -d campus-db

            Write-Host "==> [dev] Starting TypeScript DB Layer with tsx watcher..." -ForegroundColor Green
            $env:SOCKET_PATH = $SocketPath
            $dbProcess = Start-Process pnpm -ArgumentList "run dev" -WorkingDirectory "$ScriptDir\db-layer" -PassThru

            Write-Host "==> [dev] Waiting for DB layer socket to initialize at $SocketPath..." -ForegroundColor Green
            $ready = $false
            for ($i = 0; $i -lt 50; $i++) {
                if (Test-Path $SocketPath) {
                    $ready = $true
                    break
                }
                Start-Sleep -Milliseconds 200
            }

            if (-not $ready) {
                Write-Error "TypeScript DB layer socket failed to bind at $SocketPath"
                exit 1
            }
            Write-Host "==> [dev] TypeScript DB layer ready." -ForegroundColor Green

            Write-Host "==> [dev] Starting Go logical backend on :$Port..." -ForegroundColor Green
            $env:DB_SOCKET_PATH = $SocketPath
            $env:PORT = $Port
            $env:ENV = "development"
            $backendProcess = Start-Process go -ArgumentList "run ./cmd/server" -WorkingDirectory "$ScriptDir\backend" -PassThru

            Write-Host "==> [dev] Waiting for HTTP API on $BackendUrl/healthz..." -ForegroundColor Green
            for ($i = 0; $i -lt 50; $i++) {
                try {
                    $resp = Invoke-WebRequest -Uri "$BackendUrl/healthz" -UseBasicParsing -TimeoutSec 1
                    if ($resp.StatusCode -eq 200) { break }
                } catch {
                    Start-Sleep -Milliseconds 200
                }
            }
            Write-Host "==> [dev] Go backend ready." -ForegroundColor Green

            $noClient = $args -contains "--no-client"
            if ($noClient) {
                Write-Host "==> [dev] Running in headless mode (--no-client). Press Ctrl+C to stop." -ForegroundColor Yellow
                Wait-Process -Id $backendProcess.Id
            } else {
                Write-Host "==> [dev] Launching Native Fyne Client..." -ForegroundColor Green
                Push-Location "$ScriptDir\apps\client"
                $env:CAMPUS_BACKEND_URL = $BackendUrl
                $clientProcess = Start-Process go -ArgumentList "run ." -PassThru
                Pop-Location
                Wait-Process -Id $clientProcess.Id
            }
        } finally {
            Write-Host "`n========================================================================" -ForegroundColor Yellow
            Write-Host "  [dev] Initiating graceful shutdown of Campus OS services..." -ForegroundColor Yellow
            Write-Host "========================================================================" -ForegroundColor Yellow

            if ($clientProcess -and -not $clientProcess.HasExited) {
                Write-Host "==> [dev] Closing Native Client..." -ForegroundColor Yellow
                $clientProcess.CloseMainWindow() | Out-Null
                Start-Sleep -Milliseconds 300
                if (-not $clientProcess.HasExited) { Stop-Process -Id $clientProcess.Id -Force }
            }
            if ($backendProcess -and -not $backendProcess.HasExited) {
                Write-Host "==> [dev] Signaling Go backend to drain connections..." -ForegroundColor Yellow
                $backendProcess.CloseMainWindow() | Out-Null
                Start-Sleep -Seconds 1
                if (-not $backendProcess.HasExited) { Stop-Process -Id $backendProcess.Id -Force }
            }
            if ($dbProcess -and -not $dbProcess.HasExited) {
                Write-Host "==> [dev] Signaling TypeScript DB layer to stop..." -ForegroundColor Yellow
                $dbProcess.CloseMainWindow() | Out-Null
                Start-Sleep -Seconds 1
                if (-not $dbProcess.HasExited) { Stop-Process -Id $dbProcess.Id -Force }
            }
            if (Test-Path $SocketPath) { Remove-Item $SocketPath -Force }
            Write-Host "==> [dev] Clean graceful shutdown complete." -ForegroundColor Green
        }
    }
}
