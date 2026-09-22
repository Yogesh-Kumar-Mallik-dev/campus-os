#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_BOOTSTRAP_SUPERADMIN_001
# Purpose: Bootstraps the Genesis Super Admin (Chairperson) authority and exports
#          both ASCII QR docket and high-resolution printable PNG.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export SOCKET_PATH="${SOCKET_PATH:-/tmp/campus-os-dev.sock}"
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

# 1. Ensure PostgreSQL is active
if ! docker ps --format '{{.Names}}' | grep -q "^campus-postgres$"; then
  echo "==> [bootstrap] Starting Campus OS PostgreSQL service..."
  (cd "$SCRIPT_DIR" && POSTGRES_PORT="$POSTGRES_PORT" docker compose up -d campus-db)
  sleep 1
fi

# 2. Check if DB layer is already listening on SOCKET_PATH
STARTED_DB=0
if [ ! -S "$SOCKET_PATH" ]; then
  echo "==> [bootstrap] Starting temporary TypeScript DB Layer on $SOCKET_PATH..."
  rm -f "$SOCKET_PATH"
  (cd "$SCRIPT_DIR/db-layer" && SOCKET_PATH="$SOCKET_PATH" DATABASE_URL="$DATABASE_URL" pnpm run dev) &
  TEMP_DB_PID=$!
  STARTED_DB=1

  cleanup() {
    if [ "$STARTED_DB" -eq 1 ] && [ -n "${TEMP_DB_PID:-}" ]; then
      echo "==> [bootstrap] Stopping temporary DB layer (PID: $TEMP_DB_PID)..."
      kill -TERM "$TEMP_DB_PID" 2>/dev/null || true
      rm -f "$SOCKET_PATH"
    fi
  }
  trap cleanup EXIT INT TERM

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
fi

# 3. Execute bootstrap command
echo "==> [bootstrap] Executing Genesis Super Admin bootstrap..."
(cd "$SCRIPT_DIR" && DB_SOCKET_PATH="$SOCKET_PATH" go run ./backend/cmd/server bootstrap-superadmin "$@")

