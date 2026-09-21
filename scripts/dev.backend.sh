#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_BACKEND_001
# Purpose: Runs the Go Logical Backend in local development mode.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [dev:backend] Launching Go Logical Backend on :8080..."
(cd "$SCRIPT_DIR/backend" && go run ./cmd/server "$@")
