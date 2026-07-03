#!/usr/bin/env bash
# =============================================================================
# setup-postgres.sh — Buat role & database PostgreSQL untuk Learn Platform
# =============================================================================
# Default: user audi, password 090393, database learn_platform
# (sinkron dengan backend/.env dan apps/docker-compose.yml)
#
# Penggunaan:
#   make db-setup-local
#   sudo -u postgres bash scripts/setup-postgres.sh
#
# Prasyarat: PostgreSQL terinstal dan service aktif (systemctl status postgresql)
# =============================================================================

set -euo pipefail

DB_USER="${DB_USER:-audi}"
DB_PASS="${DB_PASS:-090393}"
DB_NAME="${DB_NAME:-learn_platform}"

echo "Creating role '${DB_USER}' and database '${DB_NAME}'..."

psql -v ON_ERROR_STOP=1 postgres <<SQL
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${DB_USER}') THEN
    CREATE ROLE ${DB_USER} WITH LOGIN PASSWORD '${DB_PASS}';
  ELSE
    ALTER ROLE ${DB_USER} WITH LOGIN PASSWORD '${DB_PASS}';
  END IF;
END
\$\$;
SQL

if ! psql -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" postgres | grep -q 1; then
  psql -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};" postgres
else
  echo "Database '${DB_NAME}' already exists."
fi

psql -v ON_ERROR_STOP=1 -d "${DB_NAME}" postgres <<SQL
GRANT ALL ON SCHEMA public TO ${DB_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${DB_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ${DB_USER};
SQL

echo ""
echo "Connect:"
echo "  psql -U ${DB_USER} -d postgres"
echo "  psql -U ${DB_USER} -d ${DB_NAME}"
echo ""
echo "Backend .env (folder apps/backend):"
echo "DATABASE_URL=postgres://${DB_USER}:${DB_PASS}@localhost:5432/${DB_NAME}?sslmode=disable"
