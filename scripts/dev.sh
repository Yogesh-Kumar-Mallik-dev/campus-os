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

    CLIENT_PID=""
    BACKEND_PID=""
    DB_PID=""

    CLEANED=0
    cleanup() {
      if [ "$CLEANED" -eq 1 ]; then
        return
      fi
      CLEANED=1
      trap - EXIT INT TERM
      echo -e "\n========================================================================"
      echo "  [dev] Initiating graceful shutdown of Campus OS services..."
      echo "========================================================================"

      # 1. Signal Native Client if running
      if [ -n "$CLIENT_PID" ] && kill -0 "$CLIENT_PID" 2>/dev/null; then
        echo "==> [dev] Closing Native Client (PID: $CLIENT_PID)..."
        pkill -TERM -P "$CLIENT_PID" 2>/dev/null || true
        kill -TERM "$CLIENT_PID" 2>/dev/null || true
      fi

      # 2. Signal Go Backend (triggers HTTP connection draining and server shutdown)
      if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
        echo "==> [dev] Signaling Go backend (PID: $BACKEND_PID) to drain connections..."
        pkill -TERM -P "$BACKEND_PID" 2>/dev/null || true
        kill -TERM "$BACKEND_PID" 2>/dev/null || true
      fi

      # 3. Signal TypeScript DB Layer (triggers gRPC tryShutdown and Prisma disconnect)
      if [ -n "$DB_PID" ] && kill -0 "$DB_PID" 2>/dev/null; then
        echo "==> [dev] Signaling TypeScript DB layer (PID: $DB_PID) to gracefully stop..."
        pkill -TERM -P "$DB_PID" 2>/dev/null || true
        kill -TERM "$DB_PID" 2>/dev/null || true
      fi

      # Also signal any remaining child process group jobs
      kill -TERM $(jobs -p) 2>/dev/null || true

      # 4. Wait gracefully up to 3 seconds for services to finalize
      echo "==> [dev] Waiting for background services to terminate..."
      for i in {1..15}; do
        alive=0
        if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then alive=1; fi
        if [ -n "$DB_PID" ] && kill -0 "$DB_PID" 2>/dev/null; then alive=1; fi
        if [ -n "$CLIENT_PID" ] && kill -0 "$CLIENT_PID" 2>/dev/null; then alive=1; fi
        if [ "$alive" -eq 0 ]; then
          break
        fi
        sleep 0.2
      done

      # 5. Force kill any hung process if still running after grace period
      if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
        pkill -KILL -P "$BACKEND_PID" 2>/dev/null || true
        kill -KILL "$BACKEND_PID" 2>/dev/null || true
      fi
      if [ -n "$DB_PID" ] && kill -0 "$DB_PID" 2>/dev/null; then
        pkill -KILL -P "$DB_PID" 2>/dev/null || true
        kill -KILL "$DB_PID" 2>/dev/null || true
      fi
      if [ -n "$CLIENT_PID" ] && kill -0 "$CLIENT_PID" 2>/dev/null; then
        pkill -KILL -P "$CLIENT_PID" 2>/dev/null || true
        kill -KILL "$CLIENT_PID" 2>/dev/null || true
      fi

      # 6. Unlink domain socket & reclaim port
      rm -f "$SOCKET_PATH"
      if lsof -ti:"$PORT" >/dev/null 2>&1; then
        lsof -ti:"$PORT" | xargs -r kill -9 2>/dev/null || true
      fi

      echo "========================================================================"
      echo "  [dev] Graceful shutdown complete. All resources cleanly released."
      echo "========================================================================"
      exit 0
    }
    trap cleanup EXIT INT TERM

    # Ensure port $PORT is free from any previous stale/zombie backend instances
    if lsof -ti:"$PORT" >/dev/null 2>&1; then
      echo "==> [dev] Port $PORT is occupied by a stale process. Reclaiming port..."
      lsof -ti:"$PORT" | xargs -r kill -9 2>/dev/null || true
      sleep 0.5
    fi

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
      if ! kill -0 "$DB_PID" 2>/dev/null; then
        echo "Error: TypeScript DB layer exited prematurely (PID: $DB_PID)."
        cleanup
        exit 1
      fi
      if [ -S "$SOCKET_PATH" ]; then
        break
      fi
      sleep 0.2
    done

    if [ ! -S "$SOCKET_PATH" ]; then
      echo "Error: TypeScript DB layer socket failed to bind at $SOCKET_PATH"
      cleanup
      exit 1
    fi
    echo "==> [dev] TypeScript DB layer ready."

    # 5. Start Go Logical Backend
    echo "==> [dev] Starting Go logical backend on :$PORT..."
    (cd "$SCRIPT_DIR" && DB_SOCKET_PATH="$SOCKET_PATH" PORT="$PORT" ./scripts/dev.backend.sh) &
    BACKEND_PID=$!

    # Wait for backend health
    echo "==> [dev] Waiting for HTTP API on ${CAMPUS_BACKEND_URL}/healthz..."
    BACKEND_READY=false
    for i in {1..50}; do
      if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
        echo "Error: Go logical backend exited prematurely (PID: $BACKEND_PID)."
        cleanup
        exit 1
      fi
      if curl -s "${CAMPUS_BACKEND_URL}/healthz" >/dev/null 2>&1; then
        BACKEND_READY=true
        break
      fi
      sleep 0.2
    done

    if [ "$BACKEND_READY" != true ]; then
      echo "Error: Timed out waiting for Go backend healthz endpoint at ${CAMPUS_BACKEND_URL}/healthz"
      cleanup
      exit 1
    fi
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
      wait "$BACKEND_PID" 2>/dev/null || true
      cleanup
    else
      # 6. Launch Native Client
      echo "==> [dev] Launching Native Fyne Client..."
      (cd "$SCRIPT_DIR" && CAMPUS_BACKEND_URL="$CAMPUS_BACKEND_URL" ./scripts/dev.client.sh "$@") &
      CLIENT_PID=$!

      # Wait for client to exit normally, then invoke graceful shutdown
      wait "$CLIENT_PID" 2>/dev/null || true
      cleanup
    fi
    ;;
esac
