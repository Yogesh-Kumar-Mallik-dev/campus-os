#!/usr/bin/env bash
# ==============================================================================
# Campus OS — Universal Lifecycle Orchestrator (POSIX Entrypoint)
# Usage: ./script.sh [dev|build|check|test|deps|proto:gen|envi|uenvi|flush|help]
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMMAND="${1:-help}"
shift || true

function show_help() {
  echo "Campus OS — Unified Lifecycle Orchestrator"
  echo "Usage: ./script.sh <command> [args...]"
  echo ""
  echo "Available commands:"
  echo "  dev        Start local development environment (Go backend, TS DB layer, Docker PG)"
  echo "  build      Compile Go binary, TS DB service, and SvelteKit web bundle"
  echo "  check      Run typecheckers, linters, and static analysis"
  echo "  test       Run unit and integration test suites (Go + TypeScript)"
  echo "  deps       Install Go modules and Node dependencies"
  echo "  proto:gen  Generate Go and TypeScript contracts via buf from /proto"
  echo "  envi       Start PostgreSQL, OTel, and supporting Docker Compose services"
  echo "  uenvi      Stop and remove Docker Compose development containers"
  echo "  flush      Flush local test database and state for a clean slate"
  echo "  help       Display this help message"
}

case "$COMMAND" in
  dev)
    if [ -f "$SCRIPT_DIR/scripts/dev.sh" ]; then
      "$SCRIPT_DIR/scripts/dev.sh" "$@"
    else
      echo "==> Starting Campus OS development environment..."
      docker compose up -d campus-db
    fi
    ;;
  build)
    if [ -f "$SCRIPT_DIR/scripts/build.sh" ]; then
      "$SCRIPT_DIR/scripts/build.sh" "$@"
    else
      echo "==> Building Campus OS artifacts..."
      if [ -d "$SCRIPT_DIR/backend" ]; then
        (cd "$SCRIPT_DIR/backend" && go build -v ./...)
      fi
      if [ -d "$SCRIPT_DIR/db-layer" ]; then
        (cd "$SCRIPT_DIR/db-layer" && npm run build)
      fi
    fi
    ;;
  check)
    if [ -f "$SCRIPT_DIR/scripts/check.sh" ]; then
      "$SCRIPT_DIR/scripts/check.sh" "$@"
    else
      echo "==> Running static checks and linting..."
      if [ -d "$SCRIPT_DIR/backend" ]; then
        (cd "$SCRIPT_DIR/backend" && go vet ./...)
      fi
      if [ -d "$SCRIPT_DIR/db-layer" ]; then
        (cd "$SCRIPT_DIR/db-layer" && npm run lint || true)
      fi
    fi
    ;;
  test)
    if [ -f "$SCRIPT_DIR/scripts/test.sh" ]; then
      "$SCRIPT_DIR/scripts/test.sh" "$@"
    else
      echo "==> Running automated test suites..."
      if [ -d "$SCRIPT_DIR/backend" ]; then
        (cd "$SCRIPT_DIR/backend" && go test -v -race ./...)
      fi
      if [ -d "$SCRIPT_DIR/db-layer" ]; then
        (cd "$SCRIPT_DIR/db-layer" && npm test || true)
      fi
    fi
    ;;
  deps)
    if [ -f "$SCRIPT_DIR/scripts/deps.sh" ]; then
      "$SCRIPT_DIR/scripts/deps.sh" "$@"
    else
      echo "==> Installing dependencies..."
      if [ -d "$SCRIPT_DIR/backend" ]; then
        (cd "$SCRIPT_DIR/backend" && go mod download)
      fi
      if [ -d "$SCRIPT_DIR/db-layer" ]; then
        (cd "$SCRIPT_DIR/db-layer" && npm install)
      fi
    fi
    ;;
  proto:gen)
    echo "==> Generating Protobuf stubs with buf..."
    if command -v buf >/dev/null 2>&1; then
      buf generate
    else
      echo "Error: 'buf' CLI not found. Please install buf (https://buf.build)."
      exit 1
    fi
    ;;
  envi)
    echo "==> Starting local infrastructure containers..."
    docker compose up -d
    ;;
  uenvi)
    echo "==> Tearing down local infrastructure containers..."
    docker compose down
    ;;
  flush)
    echo "==> Flushing local test database..."
    docker compose down -v
    docker compose up -d campus-db
    ;;
  help|--help|-h)
    show_help
    ;;
  *)
    echo "Unknown command: $COMMAND"
    show_help
    exit 1
    ;;
esac
