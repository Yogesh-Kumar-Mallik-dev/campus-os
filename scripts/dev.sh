#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_001
# Purpose: Launches local development environment with watchers and containers.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [dev] Starting Campus OS backend database service..."
(cd "$SCRIPT_DIR" && docker compose up -d campus-db)

echo "==> [dev] Starting TypeScript DB Layer in development mode..."
(cd "$SCRIPT_DIR/db-layer" && pnpm run dev)
