#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_CHECK_001
# Purpose: Executes static analysis, linters, and typechecking across all domains.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [check] Running Go vet and static analysis..."
(cd "$SCRIPT_DIR" && go vet ./backend/... ./apps/client/...)

echo "==> [check] Running TypeScript typecheckers..."
(cd "$SCRIPT_DIR/db-layer" && pnpm run check)

echo "==> [check] Linting Protobuf schemas..."
if command -v buf >/dev/null 2>&1; then
  (cd "$SCRIPT_DIR" && buf lint proto)
else
  (cd "$SCRIPT_DIR" && pnpm exec buf lint proto)
fi

echo "==> [check] All static analysis checks passed."
