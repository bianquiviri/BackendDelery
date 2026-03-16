#!/bin/bash

# Database Setup Script
# Waits for Postgres to be ready and verifies connectivity.

set -e

echo "⏳ Waiting for database to be ready..."

# We use the service name 'db' from docker-compose.yml
# This script should be run from the host OR inside the app container if needed.
# Since AutoMigrate runs on app startup, we just ensure DB is reachable.

MAX_RETRIES=30
RETRY_COUNT=0

until docker-compose exec -T db pg_isready -U postgres > /dev/null 2>&1 || [ $RETRY_COUNT -eq $MAX_RETRIES ]; do
  echo "...waiting for postgres ($RETRY_COUNT/$MAX_RETRIES)..."
  sleep 2
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "❌ Database timed out."
    exit 1
fi

echo "✅ Database is ready!"

# Optional: Run custom SQL if AutoMigrate isn't enough
# docker-compose exec -T db psql -U postgres -d backenddelivery -f some_seed.sql

echo "🚀 Database setup verified."
