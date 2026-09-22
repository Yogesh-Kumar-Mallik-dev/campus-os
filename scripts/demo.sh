#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEMO_001
# Purpose: Launches the Campus OS Responsive Layout Engine Demo Application.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [demo] Launching Campus OS Responsive Layout Engine Demo..."
(cd "$SCRIPT_DIR/apps/client" && go run ./cmd/demo "$@")
