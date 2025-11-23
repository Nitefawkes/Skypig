#!/bin/bash
set -e

# Database migration script for Ham-Radio Cloud

# Default configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-hamradio_cloud}"

echo "🗄️  Running database migrations..."
echo "   Host: $DB_HOST:$DB_PORT"
echo "   Database: $DB_NAME"
echo ""

# Check if PostgreSQL is available
if ! command -v psql &> /dev/null; then
    echo "❌ psql not found. Using Docker exec instead..."
    PSQL_CMD="docker exec -i hamradio_postgres psql -U $DB_USER -d $DB_NAME"
else
    PSQL_CMD="PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"
fi

# Get the directory of migrations
MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../backend/internal/database/migrations" && pwd)"

echo "📂 Migrations directory: $MIGRATIONS_DIR"
echo ""

# Run each migration in order
for migration in "$MIGRATIONS_DIR"/*.up.sql; do
    filename=$(basename "$migration")
    echo "⚙️  Applying migration: $filename"

    if command -v psql &> /dev/null; then
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
    else
        cat "$migration" | docker exec -i hamradio_postgres psql -U $DB_USER -d $DB_NAME
    fi

    if [ $? -eq 0 ]; then
        echo "   ✅ $filename applied successfully"
    else
        echo "   ❌ $filename failed"
        exit 1
    fi
    echo ""
done

echo "✅ All migrations completed successfully!"
