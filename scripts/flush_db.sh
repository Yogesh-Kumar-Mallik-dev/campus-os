#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_FLUSH_DB_001
# Purpose: Flushes local PostgreSQL volume and resets database schemas to clean slate.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [flush] Stopping and removing PostgreSQL data volume..."
(cd "$SCRIPT_DIR" && docker compose down -v)

echo "==> [flush] Re-starting clean PostgreSQL service..."
(cd "$SCRIPT_DIR" && docker compose up -d campus-db)

echo "==> [flush] Database volume flushed and fresh instance running."
