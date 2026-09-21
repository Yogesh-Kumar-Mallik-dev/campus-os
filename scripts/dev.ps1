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
        Write-Host "==> [dev] Starting Campus OS backend database service..." -ForegroundColor Green
        docker compose up -d campus-db

        Write-Host "==> [dev] Starting TypeScript DB Layer in development mode..." -ForegroundColor Green
        Push-Location "$ScriptDir\db-layer"
        try {
            pnpm run dev
        } finally {
            Pop-Location
        }
    }
}
