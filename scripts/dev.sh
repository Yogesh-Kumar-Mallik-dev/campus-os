#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_DEV_001
# Purpose: Launches local development environment with watchers and containers.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SUBCOMMAND="${1:-}"

case "$SUBCOMMAND" in
  mock)
    exec "$SCRIPT_DIR/scripts/dev.mock.sh" "${@:2}"
    ;;
  backend)
    exec "$SCRIPT_DIR/scripts/dev.backend.sh" "${@:2}"
    ;;
  db)
    exec "$SCRIPT_DIR/scripts/dev.db.sh" "${@:2}"
    ;;
  client)
    exec "$SCRIPT_DIR/scripts/dev.client.sh" "${@:2}"
    ;;
  *)
    # 1. Resolve PostgreSQL port & DATABASE_URL
    if [ -z "${POSTGRES_PORT:-}" ]; then
      RUNNING_PORT="$(docker port campus-postgres 5432 2>/dev/null | head -n1 | awk -F: '{print $NF}' || true)"
      if [ -n "$RUNNING_PORT" ]; then
        export POSTGRES_PORT="$RUNNING_PORT"
      elif ss -tulpn 2>/dev/null | grep -q ":5432 "; then
        export POSTGRES_PORT="5434"
      else
        export POSTGRES_PORT="5432"
      fi
    fi
    export DATABASE_URL="${DATABASE_URL:-postgresql://campus_admin:campus_password_dev@localhost:${POSTGRES_PORT}/campus_os?schema=auth_schema}"

    # 2. Defaults for socket, backend port, and client
    export SOCKET_PATH="${SOCKET_PATH:-/tmp/campus-os-dev.sock}"
    export PORT="${PORT:-8080}"
    export CAMPUS_BACKEND_URL="${CAMPUS_BACKEND_URL:-http://localhost:$PORT}"

    echo "========================================================================"
    echo "  CAMPUS OS: LOCAL FULL-STACK DEVELOPMENT ENVIRONMENT                  "
    echo "========================================================================"
    echo "  PostgreSQL Port:    ${POSTGRES_PORT}"
    echo "  Persistence Socket: unix://${SOCKET_PATH}"
    echo "  Backend API:        ${CAMPUS_BACKEND_URL}"
    echo "  Client UI:          Native Fyne (Go)"
    echo "========================================================================"

    # Cleanup background processes on exit
    cleanup() {
      echo -e "\n==> [dev] Shutting down development services..."
      kill $(jobs -p) 2>/dev/null || true
      rm -f "$SOCKET_PATH"
      echo "==> [dev] Clean shutdown complete."
    }
    trap cleanup EXIT INT TERM

    # 3. Start PostgreSQL Container
    echo "==> [dev] Ensuring Campus OS PostgreSQL database is running (Port: $POSTGRES_PORT)..."
    (cd "$SCRIPT_DIR" && POSTGRES_PORT="$POSTGRES_PORT" docker compose up -d campus-db)

    # 4. Start TypeScript DB Layer
    rm -f "$SOCKET_PATH"
    echo "==> [dev] Starting TypeScript DB Layer with tsx watcher..."
    (cd "$SCRIPT_DIR" && SOCKET_PATH="$SOCKET_PATH" DATABASE_URL="$DATABASE_URL" ./scripts/dev.db.sh) &
    DB_PID=$!

    # Wait for socket to become ready
    echo "==> [dev] Waiting for DB layer socket to initialize at $SOCKET_PATH..."
    for i in {1..50}; do
      if [ -S "$SOCKET_PATH" ]; then
        break
      fi
      sleep 0.2
    done

    if [ ! -S "$SOCKET_PATH" ]; then
      echo "Error: TypeScript DB layer socket failed to bind at $SOCKET_PATH"
      exit 1
    fi
    echo "==> [dev] TypeScript DB layer ready."

    # 5. Start Go Logical Backend
    echo "==> [dev] Starting Go logical backend on :$PORT..."
    (cd "$SCRIPT_DIR" && DB_SOCKET_PATH="$SOCKET_PATH" PORT="$PORT" ./scripts/dev.backend.sh) &
    BACKEND_PID=$!

    # Wait for backend health
    echo "==> [dev] Waiting for HTTP API on ${CAMPUS_BACKEND_URL}/healthz..."
    for i in {1..50}; do
      if curl -s "${CAMPUS_BACKEND_URL}/healthz" >/dev/null 2>&1; then
        break
      fi
      sleep 0.2
    done
    echo "==> [dev] Go backend ready."

    # Check for headless flag
    NO_CLIENT=false
    for arg in "$@"; do
      if [ "$arg" == "--no-client" ]; then
        NO_CLIENT=true
      fi
    done

    if [ "$NO_CLIENT" = true ]; then
      echo "==> [dev] Running in headless mode (--no-client). Press Ctrl+C to stop."
      wait
    else
      # 6. Launch Native Client
      echo "==> [dev] Launching Native Fyne Client..."
      CAMPUS_BACKEND_URL="$CAMPUS_BACKEND_URL" "$SCRIPT_DIR/scripts/dev.client.sh" "$@" || true
      wait
    fi
    ;;
esac
