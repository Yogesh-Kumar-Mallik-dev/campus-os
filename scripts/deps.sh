#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEPS_001
# Purpose: Installs dependencies across Go modules and pnpm workspace packages.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [deps] Verifying prerequisites..."
command -v go >/dev/null 2>&1 || { echo "Error: 'go' is not installed."; exit 1; }
command -v pnpm >/dev/null 2>&1 || { echo "Error: 'pnpm' is not installed."; exit 1; }

echo "==> [deps] Installing pnpm workspace dependencies..."
(cd "$SCRIPT_DIR" && pnpm install)

echo "==> [deps] Downloading Go module dependencies..."
(cd "$SCRIPT_DIR/backend" && go mod download)
(cd "$SCRIPT_DIR/apps/client" && go mod download)

echo "==> [deps] Dependencies successfully installed across all workspaces."
