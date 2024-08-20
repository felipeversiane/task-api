#!/bin/bash
set -e

until pg_isready -q -d "$POSTGRES_DB" -U "$POSTGRES_USER"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "Configuring PostgreSQL settings..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    ALTER SYSTEM SET max_connections = 300;
    ALTER SYSTEM SET shared_buffers TO '425MB';
EOSQL

echo "Executing SQL commands..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS pg_trgm;
    ALTER DATABASE postgres SET synchronous_commit = OFF;
EOSQL

echo "PostgreSQL configuration and SQL commands executed successfully."
