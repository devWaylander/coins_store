-- migrate:up
CREATE SCHEMA shop;

-- balance
CREATE TABLE shop."balance" (
    id BIGSERIAL PRIMARY KEY,
    amount BIGINT NOT NULL DEFAULT 1000 CHECK (amount >= 0),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- balance history
CREATE TABLE shop."balance_history" (
    id BIGSERIAL PRIMARY KEY,
    balance_id BIGINT NOT NULL REFERENCES shop."balance" (id),
    transaction_amount BIGINT NOT NULL,
    sender VARCHAR(64) NOT NULL,
    recipient VARCHAR(64) NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "balance_history@balance_id_idx" ON shop."balance_history" (balance_id);

-- user
CREATE TABLE shop."user" (
    id BIGSERIAL PRIMARY KEY,
    balance_id BIGINT NOT NULL REFERENCES shop."balance" (id),
    username VARCHAR(64) NOT NULL,
    password_hash CHAR(64) DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX "user@username_idx" ON shop."user" (username);
CREATE INDEX "user@balance_id_idx" ON shop."user" (balance_id);

-- merch
CREATE TABLE shop."merch" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "merch@name_id_idx" ON shop."merch" (name);

-- Исходные данные
INSERT INTO shop."merch" (name, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);

-- inventory
CREATE TABLE shop."inventory" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES shop."user" (id),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- inventory_merch
CREATE TABLE shop."inventory_merch" (
    PRIMARY KEY (inventory_id, merch_id),
    inventory_id BIGINT NOT NULL REFERENCES shop."inventory" (id),
    merch_id BIGINT NOT NULL REFERENCES shop."merch" (id),
    name VARCHAR(255) UNIQUE NOT NULL,
    count BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "inventory_merch@inventory_id_idx" ON shop."inventory_merch" (inventory_id);
CREATE INDEX "inventory_merch@merch_id_idx" ON shop."inventory_merch" (merch_id);

-- migrate:down
DROP SCHEMA IF EXISTS shop CASCADE;