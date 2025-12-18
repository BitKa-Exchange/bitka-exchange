#!/bin/bash

set -e
set -u

# Function to create a user and database safely
# Usage: create_user_and_db <database_name>
create_user_and_db() {
    local database=$1
    echo "  Creating user and database '$database'"
    
    # We create a USER with the same name as the DB for simplicity/security isolation
    # We check if database exists first to prevent errors if script is run manually
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        -- Create User if not exists (Hack: Postgres doesn't support IF NOT EXISTS for users cleanly in SQL)
        DO
        \$do\$
        BEGIN
           IF NOT EXISTS (
              SELECT FROM pg_catalog.pg_roles
              WHERE  rolname = '$database') THEN
              CREATE ROLE $database LOGIN PASSWORD 'password'; -- Set default password or use env var
           END IF;
        END
        \$do\$;

        -- Create Database
        SELECT 'CREATE DATABASE $database OWNER $database'
        WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$database')\gexec
EOSQL
}

# Main Execution Logic
if [ -n "${POSTGRES_MULTIPLE_DATABASES:-}" ]; then
    echo "🚀 Initializing multiple databases: $POSTGRES_MULTIPLE_DATABASES"
    
    # Convert commas to spaces to loop
    for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
        create_user_and_db $db
    done
    
    echo "✅ All requested databases created!"
else
    echo "⚠️ No additional databases specified in POSTGRES_MULTIPLE_DATABASES"
fi