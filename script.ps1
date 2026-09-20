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
        Write-Host "==> Starting development environment..." -ForegroundColor Green
        docker compose up -d campus-db
    }
    "build" {
        Write-Host "==> Building Campus OS artifacts..." -ForegroundColor Green
        if (Test-Path "$ScriptDir\backend") {
            Push-Location "$ScriptDir\backend"; go build -v ./...; Pop-Location
        }
        if (Test-Path "$ScriptDir\db-layer") {
            Push-Location "$ScriptDir\db-layer"; npm run build; Pop-Location
        }
    }
    "check" {
        Write-Host "==> Running static analysis..." -ForegroundColor Green
        if (Test-Path "$ScriptDir\backend") {
            Push-Location "$ScriptDir\backend"; go vet ./...; Pop-Location
        }
        if (Test-Path "$ScriptDir\db-layer") {
            Push-Location "$ScriptDir\db-layer"; npm run lint; Pop-Location
        }
    }
    "test" {
        Write-Host "==> Running test suites..." -ForegroundColor Green
        if (Test-Path "$ScriptDir\backend") {
            Push-Location "$ScriptDir\backend"; go test -v -race ./...; Pop-Location
        }
        if (Test-Path "$ScriptDir\db-layer") {
            Push-Location "$ScriptDir\db-layer"; npm test; Pop-Location
        }
    }
    "deps" {
        Write-Host "==> Installing dependencies..." -ForegroundColor Green
        if (Test-Path "$ScriptDir\backend") {
            Push-Location "$ScriptDir\backend"; go mod download; Pop-Location
        }
        if (Test-Path "$ScriptDir\db-layer") {
            Push-Location "$ScriptDir\db-layer"; npm install; Pop-Location
        }
    }
    "proto:gen" {
        Write-Host "==> Generating Protobuf stubs..." -ForegroundColor Green
        buf generate
    }
    "envi" {
        Write-Host "==> Starting Docker Compose services..." -ForegroundColor Green
        docker compose up -d
    }
    "uenvi" {
        Write-Host "==> Tearing down Docker Compose services..." -ForegroundColor Green
        docker compose down
    }
    "flush" {
        Write-Host "==> Flushing local test database..." -ForegroundColor Yellow
        docker compose down -v
        docker compose up -d campus-db
    }
    Default {
        Show-Help
    }
}
