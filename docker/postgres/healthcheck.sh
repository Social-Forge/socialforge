#!/bin/sh
set -e

# Test PostgreSQL connection
pg_isready -U "${POSTGRES_USER:-forge}" -d "${POSTGRES_DB:-forge}"

# Test jika extensions tersedia
psql -U "${POSTGRES_USER:-forge}" -d "${POSTGRES_DB:-forge}" -c "SELECT version();" > /dev/null

echo "PostgreSQL is healthy"
