# ==============================================================================
# Campus OS — Universal Lifecycle Orchestrator (Windows PowerShell Entrypoint)
# Usage: .\script.ps1 [dev|build|check|test|deps|proto:gen|envi|uenvi|flush|help]
# ==============================================================================
[CmdletBinding()]
param (
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Show-Help {
    Write-Host "Campus OS — Unified Lifecycle Orchestrator" -ForegroundColor Cyan
    Write-Host "Usage: .\script.ps1 <command>"
    Write-Host ""
    Write-Host "Available commands:"
    Write-Host "  dev        Start local development environment"
    Write-Host "  build      Compile Go binary, TS DB service, and web bundle"
    Write-Host "  check      Run typecheckers, linters, and static analysis"
    Write-Host "  test       Run unit and integration test suites"
    Write-Host "  deps       Install Go modules and Node dependencies"
    Write-Host "  proto:gen  Generate Go and TypeScript contracts via buf"
    Write-Host "  envi       Start PostgreSQL and Docker Compose services"
    Write-Host "  uenvi      Stop Docker Compose containers"
    Write-Host "  flush      Flush database for clean slate"
    Write-Host "  help       Display this help message"
}

switch ($Command) {
    "dev" {
        if (Test-Path "$ScriptDir\scripts\dev.ps1") { & "$ScriptDir\scripts\dev.ps1" @args }
        else { docker compose up -d campus-db; Push-Location "$ScriptDir\db-layer"; pnpm run dev; Pop-Location }
    }
    "build" {
        if (Test-Path "$ScriptDir\scripts\build.ps1") { & "$ScriptDir\scripts\build.ps1" @args }
        else {
            Push-Location "$ScriptDir\backend"; go build -v ./...; Pop-Location
            Push-Location "$ScriptDir\db-layer"; pnpm run build; Pop-Location
        }
    }
    "check" {
        if (Test-Path "$ScriptDir\scripts\check.ps1") { & "$ScriptDir\scripts\check.ps1" @args }
        else {
            Push-Location "$ScriptDir\backend"; go vet ./...; Pop-Location
            Push-Location "$ScriptDir\db-layer"; pnpm run check; Pop-Location
        }
    }
    "test" {
        if (Test-Path "$ScriptDir\scripts\test.ps1") { & "$ScriptDir\scripts\test.ps1" @args }
        else {
            go test -v -race ./backend/... ./apps/client/...
            Push-Location "$ScriptDir\db-layer"; pnpm test; Pop-Location
        }
    }
    "deps" {
        if (Test-Path "$ScriptDir\scripts\deps.ps1") { & "$ScriptDir\scripts\deps.ps1" @args }
        else {
            pnpm install
            Push-Location "$ScriptDir\backend"; go mod download; Pop-Location
        }
    }
    "proto:gen" {
        if (Test-Path "$ScriptDir\scripts\proto_gen.ps1") { & "$ScriptDir\scripts\proto_gen.ps1" @args }
        else { pnpm exec buf generate proto --template proto/buf.gen.yaml }
    }
    "envi" {
        if (Test-Path "$ScriptDir\scripts\envi.ps1") { & "$ScriptDir\scripts\envi.ps1" @args }
        else { docker compose up -d campus-db }
    }
    "uenvi" {
        if (Test-Path "$ScriptDir\scripts\uenvi.ps1") { & "$ScriptDir\scripts\uenvi.ps1" @args }
        else { docker compose down }
    }
    "flush" {
        if (Test-Path "$ScriptDir\scripts\flush_db.ps1") { & "$ScriptDir\scripts\flush_db.ps1" @args }
        else { docker compose down -v; docker compose up -d campus-db }
    }
    Default {
        Show-Help
    }
}
