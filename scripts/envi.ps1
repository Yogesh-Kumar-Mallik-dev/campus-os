# BLOCK_SCRIPT_ENVI_PS_001
$ErrorActionPreference = "Stop"

Write-Host "==> [envi] Starting Campus OS local infrastructure..." -ForegroundColor Green
docker compose up -d campus-db
