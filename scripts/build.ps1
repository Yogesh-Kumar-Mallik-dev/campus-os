# BLOCK_SCRIPT_BUILD_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [build] Building Go Logical Backend..." -ForegroundColor Green
New-Item -ItemType Directory -Force -Path "$ScriptDir\bin" | Out-Null
Push-Location "$ScriptDir\backend"; go build -v -o "$ScriptDir\bin\campus-backend.exe" ./cmd/server; Pop-Location

Write-Host "==> [build] Building TypeScript DB Layer..." -ForegroundColor Green
Push-Location "$ScriptDir\db-layer"; pnpm run build; Pop-Location

Write-Host "==> [build] All production artifacts built successfully." -ForegroundColor Green
