# BLOCK_SCRIPT_CHECK_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [check] Running Go vet and static analysis..." -ForegroundColor Green
go vet ./backend/... ./apps/client/...

Write-Host "==> [check] Running TypeScript typecheckers..." -ForegroundColor Green
Push-Location "$ScriptDir\db-layer"; pnpm run check; Pop-Location

Write-Host "==> [check] Linting Protobuf schemas..." -ForegroundColor Green
pnpm exec buf lint proto

Write-Host "==> [check] All static analysis checks passed." -ForegroundColor Green
