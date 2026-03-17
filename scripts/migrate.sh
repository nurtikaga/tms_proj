#!/usr/bin/env bash
set -euo pipefail

MIGRATIONS_PATH="${1:-./scripts/migrations}"
DATABASE_URL="${POSTGRES_DSN:-postgres://tms:tms@localhost:5432/tms?sslmode=disable}"

echo "Running migrations from: $MIGRATIONS_PATH"
echo "Target database:         $DATABASE_URL"

migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up

echo "Migrations applied successfully."
