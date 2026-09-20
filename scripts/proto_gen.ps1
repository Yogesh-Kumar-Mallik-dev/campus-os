# BLOCK_SCRIPT_PROTO_GEN_PS_001
$ErrorActionPreference = "Stop"

Write-Host "==> [proto:gen] Linting Protobuf schemas..." -ForegroundColor Green
pnpm exec buf lint proto

Write-Host "==> [proto:gen] Generating language contracts..." -ForegroundColor Green
pnpm exec buf generate proto --template proto/buf.gen.yaml

Write-Host "==> [proto:gen] Protobuf contracts generated successfully." -ForegroundColor Green
