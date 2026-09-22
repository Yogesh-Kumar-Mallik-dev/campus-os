#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_BACKEND_001
# Purpose: Runs the Go Logical Backend in local development mode.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export DB_SOCKET_PATH="${DB_SOCKET_PATH:-/tmp/campus-os-dev.sock}"
export PORT="${PORT:-8080}"
export ENV="${ENV:-development}"

echo "==> [dev:backend] Launching Go Logical Backend on :$PORT (Socket: $DB_SOCKET_PATH)..."
cd "$SCRIPT_DIR/backend"
exec go run ./cmd/server "$@"
