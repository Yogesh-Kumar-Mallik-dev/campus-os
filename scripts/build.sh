#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_BUILD_001
# Purpose: Compiles production assets, binaries, and packages.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [build] Building Go Logical Backend..."
(cd "$SCRIPT_DIR/backend" && go build -v -o "$SCRIPT_DIR/bin/campus-backend" ./cmd/server)

echo "==> [build] Building TypeScript DB Layer..."
(cd "$SCRIPT_DIR/db-layer" && pnpm run build)

echo "==> [build] Building Native Go Client (Fyne)..."
(cd "$SCRIPT_DIR/apps/client" && go build -v -o "$SCRIPT_DIR/bin/campus-client" .)

echo "==> [build] All production artifacts built successfully."
