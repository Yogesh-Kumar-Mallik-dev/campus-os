#!/usr/bin/env bash
# ==============================================================================
# BLOCK_SCRIPT_PROTO_GEN_001
# Purpose: Compiles Protobuf definitions using buf into Go and TypeScript stubs.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [proto:gen] Checking buf compiler..."
if command -v buf >/dev/null 2>&1; then
  BUF_CMD="buf"
elif command -v pnpm >/dev/null 2>&1; then
  BUF_CMD="pnpm exec buf"
else
  echo "Error: Neither 'buf' nor 'pnpm' found."
  exit 1
fi

echo "==> [proto:gen] Linting Protobuf schemas..."
(cd "$SCRIPT_DIR" && $BUF_CMD lint proto)

echo "==> [proto:gen] Generating language contracts..."
(cd "$SCRIPT_DIR" && $BUF_CMD generate proto --template proto/buf.gen.yaml)

echo "==> [proto:gen] Protobuf contracts generated successfully."
