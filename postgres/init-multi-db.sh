#!/bin/bash
set -e

# 1. Create the databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE items_db;
	CREATE DATABASE orders_db;
EOSQL

# 2. Initialize items_db schema and seed data
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "items_db" <<-EOSQL
    CREATE TABLE IF NOT EXISTS items (
        id VARCHAR(50) PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        price DECIMAL(10,2) NOT NULL,
        stock INT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP WITH TIME ZONE
    );

    CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items(deleted_at);

    INSERT INTO items (id, title, price, stock) 
    VALUES ('MLA43960787', 'Monitor gamer curvo Xiaomi Gaming G34WQi LCD negro', 619999.00, 55)
    ON CONFLICT (id) DO NOTHING;
EOSQL

# 3. Initialize orders_db schema
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "orders_db" <<-EOSQL
    CREATE TABLE IF NOT EXISTS orders (
        id VARCHAR(50) PRIMARY KEY,
        user_id VARCHAR(100) NOT NULL,
        item_id VARCHAR(50) NOT NULL,
        quantity INT NOT NULL,
        amount NUMERIC(12,2) NOT NULL,
        address VARCHAR(255) NOT NULL,
        status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP WITH TIME ZONE
    );

    CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
EOSQL
