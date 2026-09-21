# BLOCK_SCRIPT_UENVI_PS_001
$ErrorActionPreference = "Stop"

Write-Host "==> [uenvi] Gracefully stopping Campus OS containers..." -ForegroundColor Green
docker compose down --timeout 10
