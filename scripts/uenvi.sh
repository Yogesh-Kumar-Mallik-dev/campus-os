#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_UENVI_001
# Purpose: Stops local Docker Compose infrastructure and cleans up runtime state.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [uenvi] Tearing down Campus OS containers..."
(cd "$SCRIPT_DIR" && docker compose down)

echo "==> [uenvi] Environment successfully terminated."
