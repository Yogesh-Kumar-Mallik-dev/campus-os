#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_DB_001
# Purpose: Runs the TypeScript DB Layer in development watch mode.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [dev:db] Launching TypeScript DB Layer with tsx watcher..."
(cd "$SCRIPT_DIR/db-layer" && pnpm run dev)
