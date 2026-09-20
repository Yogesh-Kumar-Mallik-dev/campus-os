#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_ENVI_001
# Purpose: Initializes local environment, network volumes, and Docker Compose.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [envi] Checking Docker engine..."
command -v docker >/dev/null 2>&1 || { echo "Error: 'docker' is not installed or running."; exit 1; }

echo "==> [envi] Starting Campus OS local infrastructure (PostgreSQL & UDS Socket)..."
(cd "$SCRIPT_DIR" && docker compose up -d campus-db)

echo "==> [envi] Environment successfully started and ready for development."
