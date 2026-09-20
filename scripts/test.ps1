# BLOCK_SCRIPT_TEST_PS_001
$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Write-Host "==> [test] Running Go test suites..." -ForegroundColor Green
go test -v -race ./backend/... ./apps/client/...

Write-Host "==> [test] Running TypeScript test suites..." -ForegroundColor Green
Push-Location "$ScriptDir\db-layer"; pnpm test; Pop-Location

Write-Host "==> [test] All tests completed successfully." -ForegroundColor Green
