#!/usr/bin/env bash
set -euo pipefail

# Run a one-off seed against the configured database.
# Usage: DATABASE_URL=postgres://user:pass@host:5432/db JWT_SECRET=... bash scripts/seed.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

DB_URL="${DATABASE_URL:-}"
if [[ -z "${DB_URL}" ]]; then
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_USER="${DB_USER:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-postgres}"
  DB_NAME="${DB_NAME:-bsmart}"
  DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
  echo "INFO: DATABASE_URL not set, using built URL from env vars (host=${DB_HOST}, db=${DB_NAME}, user=${DB_USER})"
else
  echo "INFO: Using provided DATABASE_URL (details not shown)"
fi

APP_ENV="${APP_ENV:-production}"
HTTP_PORT="${HTTP_PORT:-8080}"
JWT_SECRET="${JWT_SECRET:-dev-secret}"
JWT_EXPIRATION="${JWT_EXPIRATION:-1h}"
WS_ALLOWED_ORIGINS="${WS_ALLOWED_ORIGINS:-http://localhost:8080,https://ignimbrite.github.io,https://8113c6fc74a6.ngrok-free.app}"

echo "Starting seed with APP_ENV=${APP_ENV} SEED_ON_START=true"

SEED_ON_START=true \
APP_ENV="${APP_ENV}" \
HTTP_PORT="${HTTP_PORT}" \
DATABASE_URL="${DB_URL}" \
JWT_SECRET="${JWT_SECRET}" \
JWT_EXPIRATION="${JWT_EXPIRATION}" \
WS_ALLOWED_ORIGINS="${WS_ALLOWED_ORIGINS}" \
  go run "${ROOT_DIR}/cmd/api"
