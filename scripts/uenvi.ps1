# BLOCK_SCRIPT_UENVI_PS_001
$ErrorActionPreference = "Stop"

Write-Host "==> [uenvi] Tearing down Campus OS containers..." -ForegroundColor Green
docker compose down
