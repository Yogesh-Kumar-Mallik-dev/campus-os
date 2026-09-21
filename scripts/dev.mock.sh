#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_MOCK_001
# Purpose: Launches zero-dependency full-stack development mode using in-memory mock persistence.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOCKET_PATH="/tmp/campus-os-mock.sock"

echo "========================================================================"
echo "  CAMPUS OS: ZERO-DEPENDENCY MOCK DEVELOPMENT ENVIRONMENT               "
echo "========================================================================"
echo "  Mock Socket Path: unix://${SOCKET_PATH}"
echo "  Backend API:      http://localhost:8080"
echo "========================================================================"

# Cleanup background processes on exit
cleanup() {
  echo -e "\n==> [dev:mock] Shutting down mock development services..."
  kill $(jobs -p) 2>/dev/null || true
  rm -f "$SOCKET_PATH"
  echo "==> [dev:mock] Clean shutdown complete."
}
trap cleanup EXIT INT TERM

# 1. Start In-Memory Mock DB Layer
echo "==> [dev:mock] Starting in-memory mock persistence server..."
SOCKET_PATH="$SOCKET_PATH" pnpm --filter @campus-os/db-layer exec tsx src/server.mock.ts &
DB_PID=$!

# Wait for socket to become ready
echo "==> [dev:mock] Waiting for mock socket to initialize..."
for i in {1..30}; do
  if [ -S "$SOCKET_PATH" ]; then
    break
  fi
  sleep 0.2
done

if [ ! -S "$SOCKET_PATH" ]; then
  echo "Error: Mock persistence server socket failed to bind at $SOCKET_PATH"
  exit 1
fi
echo "==> [dev:mock] Mock persistence server ready."

# 2. Start Go Logical Backend
echo "==> [dev:mock] Starting Go logical backend..."
(cd "$SCRIPT_DIR/backend" && DB_SOCKET_PATH="$SOCKET_PATH" PORT="8080" go run ./cmd/server) &
BACKEND_PID=$!

# Wait for backend health
echo "==> [dev:mock] Waiting for HTTP API on http://localhost:8080/healthz..."
for i in {1..30}; do
  if curl -s http://localhost:8080/healthz >/dev/null 2>&1; then
    break
  fi
  sleep 0.2
done
echo "==> [dev:mock] Go backend ready."

# 3. Launch Native Client
echo "==> [dev:mock] Launching Native Fyne Client..."
(cd "$SCRIPT_DIR/apps/client" && go run .) || true

wait
