# BLOCK_SCRIPT_FLUSH_DB_PS_001
$ErrorActionPreference = "Stop"

Write-Host "==> [flush] Flushing database data..." -ForegroundColor Yellow
docker compose down -v
docker compose up -d campus-db
