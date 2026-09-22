# ==============================================================================
# BLOCK_SCRIPT_BOOTSTRAP_SUPERADMIN_PS1_001
# Purpose: Bootstraps Genesis Super Admin authority via Prisma DB layer over UDS.
# ==============================================================================
$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$SocketPath = if ($env:SOCKET_PATH) { $env:SOCKET_PATH } else { "/tmp/campus-os-dev.sock" }
$PostgresPort = if ($env:POSTGRES_PORT) { $env:POSTGRES_PORT } else { "5434" }
$DatabaseUrl = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { "postgresql://campus_admin:campus_secure_password_2026@127.0.0.1:$PostgresPort/campus_os?schema=auth_schema" }

Write-Host "==> [bootstrap] Executing Genesis Super Admin bootstrap via Prisma DB layer..."
Push-Location "$ScriptDir/backend"
try {
    $env:DB_SOCKET_PATH = $SocketPath
    go run ./cmd/server bootstrap-superadmin $args
} finally {
    Pop-Location
}
