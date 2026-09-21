#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_CLIENT_001
# Purpose: Runs the Go Fyne Native Client in development mode.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [dev:client] Launching Go Fyne Native Client..."
(cd "$SCRIPT_DIR/apps/client" && go run .)
