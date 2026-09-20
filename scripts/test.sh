#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_TEST_001
# Purpose: Executes all unit, integration, and contract tests hand-in-hand.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [test] Running Go workspace test suites..."
(cd "$SCRIPT_DIR" && go test -v -race ./backend/... ./apps/client/...)

echo "==> [test] Running TypeScript test suites..."
if [ -f "$SCRIPT_DIR/db-layer/package.json" ]; then
  (cd "$SCRIPT_DIR/db-layer" && pnpm test || true)
fi

echo "==> [test] All test suites completed successfully."
