#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_DB_001
# Purpose: Runs the TypeScript DB Layer in development watch mode.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export SOCKET_PATH="${SOCKET_PATH:-/tmp/campus-os-dev.sock}"

echo "==> [dev:db] Launching TypeScript DB Layer with tsx watcher on $SOCKET_PATH..."
(cd "$SCRIPT_DIR/db-layer" && SOCKET_PATH="$SOCKET_PATH" pnpm run dev "$@")
