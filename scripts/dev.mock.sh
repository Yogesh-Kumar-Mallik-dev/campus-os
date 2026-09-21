#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_MOCK_001
# Purpose: Launches zero-dependency full-stack development mode using in-memory mock persistence.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOCKET_PATH="${SOCKET_PATH:-/tmp/campus-os-mock.sock}"
PORT="${PORT:-8080}"
CAMPUS_BACKEND_URL="${CAMPUS_BACKEND_URL:-http://localhost:$PORT}"

echo "========================================================================"
echo "  CAMPUS OS: ZERO-DEPENDENCY MOCK DEVELOPMENT ENVIRONMENT               "
echo "========================================================================"
echo "  Mock Socket Path: unix://${SOCKET_PATH}"
echo "  Backend API:      ${CAMPUS_BACKEND_URL}"
echo "========================================================================"

CLIENT_PID=""
BACKEND_PID=""
DB_PID=""

# Graceful shutdown handler
cleanup() {
  trap - EXIT INT TERM
  echo -e "\n========================================================================"
  echo "  [dev:mock] Initiating graceful shutdown of mock services..."
  echo "========================================================================"

  if [ -n "$CLIENT_PID" ] && kill -0 "$CLIENT_PID" 2>/dev/null; then
    echo "==> [dev:mock] Closing Native Client (PID: $CLIENT_PID)..."
    kill -TERM "$CLIENT_PID" 2>/dev/null || true
  fi

  if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
    echo "==> [dev:mock] Signaling Go backend (PID: $BACKEND_PID) to drain connections..."
    kill -TERM "$BACKEND_PID" 2>/dev/null || true
  fi

  if [ -n "$DB_PID" ] && kill -0 "$DB_PID" 2>/dev/null; then
    echo "==> [dev:mock] Signaling Mock DB layer (PID: $DB_PID) to stop..."
    kill -TERM "$DB_PID" 2>/dev/null || true
  fi

  kill -TERM $(jobs -p) 2>/dev/null || true

  echo "==> [dev:mock] Waiting for background services to terminate..."
  for i in {1..25}; do
    alive=0
    if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then alive=1; fi
    if [ -n "$DB_PID" ] && kill -0 "$DB_PID" 2>/dev/null; then alive=1; fi
    if [ -n "$CLIENT_PID" ] && kill -0 "$CLIENT_PID" 2>/dev/null; then alive=1; fi
    if [ "$alive" -eq 0 ]; then
      break
    fi
    sleep 0.2
  done

  if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
    kill -KILL "$BACKEND_PID" 2>/dev/null || true
  fi
  if [ -n "$DB_PID" ] && kill -0 "$DB_PID" 2>/dev/null; then
    kill -KILL "$DB_PID" 2>/dev/null || true
  fi
  if [ -n "$CLIENT_PID" ] && kill -0 "$CLIENT_PID" 2>/dev/null; then
    kill -KILL "$CLIENT_PID" 2>/dev/null || true
  fi

  rm -f "$SOCKET_PATH"
  echo "========================================================================"
  echo "  [dev:mock] Clean graceful shutdown complete."
  echo "========================================================================"
}
trap cleanup EXIT INT TERM

# 1. Start In-Memory Mock DB Layer
rm -f "$SOCKET_PATH"
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
(cd "$SCRIPT_DIR/backend" && DB_SOCKET_PATH="$SOCKET_PATH" PORT="$PORT" go run ./cmd/server) &
BACKEND_PID=$!

# Wait for backend health
echo "==> [dev:mock] Waiting for HTTP API on ${CAMPUS_BACKEND_URL}/healthz..."
for i in {1..30}; do
  if curl -s "${CAMPUS_BACKEND_URL}/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 0.2
done
echo "==> [dev:mock] Go backend ready."

# Check for headless flag
NO_CLIENT=false
for arg in "$@"; do
  if [ "$arg" == "--no-client" ]; then
    NO_CLIENT=true
  fi
done

if [ "$NO_CLIENT" = true ]; then
  echo "==> [dev:mock] Running in headless mode (--no-client). Press Ctrl+C to stop."
  wait
else
  # 3. Launch Native Client
  echo "==> [dev:mock] Launching Native Fyne Client..."
  (cd "$SCRIPT_DIR/apps/client" && CAMPUS_BACKEND_URL="$CAMPUS_BACKEND_URL" go run .) &
  CLIENT_PID=$!
  wait "$CLIENT_PID" 2>/dev/null || true
  cleanup
fi
