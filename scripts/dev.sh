#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_001
# Purpose: Launches local development environment with watchers and containers.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SUBCOMMAND="${1:-}"

case "$SUBCOMMAND" in
  mock)
    exec "$SCRIPT_DIR/scripts/dev.mock.sh" "${@:2}"
    ;;
  backend)
    exec "$SCRIPT_DIR/scripts/dev.backend.sh" "${@:2}"
    ;;
  db)
    exec "$SCRIPT_DIR/scripts/dev.db.sh" "${@:2}"
    ;;
  client)
    exec "$SCRIPT_DIR/scripts/dev.client.sh" "${@:2}"
    ;;
  *)
    echo "==> [dev] Starting Campus OS backend database service..."
    (cd "$SCRIPT_DIR" && docker compose up -d campus-db)

    echo "==> [dev] Starting TypeScript DB Layer in development mode..."
    (cd "$SCRIPT_DIR/db-layer" && pnpm run dev)
    ;;
esac
