#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_CLIENT_001
# Purpose: Runs the Go Fyne Native Client in development mode.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export CAMPUS_BACKEND_URL="${CAMPUS_BACKEND_URL:-http://localhost:8080}"

echo "==> [dev:client] Launching Go Fyne Native Client (Backend: $CAMPUS_BACKEND_URL)..."
(cd "$SCRIPT_DIR/apps/client" && CAMPUS_BACKEND_URL="$CAMPUS_BACKEND_URL" go run . "$@")
