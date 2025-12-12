#!/bin/bash
set -e

# Function to create a user and database safely
create_db() {
    local db=$1
    echo "  Creating user and database '$db'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        CREATE USER $db;
        CREATE DATABASE $db;
        GRANT ALL PRIVILEGES ON DATABASE $db TO $db;
EOSQL
}

# Add your microservice databases here
# We use the generic 'POSTGRES_USER' (from .env) as the superuser to create these
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE bitka_auth;
    CREATE DATABASE bitka_account;
    CREATE DATABASE bitka_ledger;
    CREATE DATABASE bitka_order;
    -- Add more as you grow
EOSQL

echo "✅ Multiple databases created successfully"