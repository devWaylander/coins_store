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

-- merch
CREATE TABLE shop."merch" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "merch@name_id_idx" ON shop."merch" (name);

-- inventory
CREATE TABLE shop."inventory" (
    id BIGSERIAL PRIMARY KEY,
    merch_id BIGINT NOT NULL REFERENCES shop."merch" (id),
    count BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "inventory@merch_id_idx" ON shop."inventory" (merch_id);

-- user
CREATE TABLE shop."user" (
    id BIGSERIAL PRIMARY KEY,
    balance_id BIGINT NOT NULL REFERENCES shop."balance" (id),
    inventory_id BIGINT DEFAULT NULL REFERENCES shop."inventory" (id),
    username VARCHAR(64) NOT NULL,
    password_hash CHAR(64) DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX "user@username_idx" ON shop."user" (username);
CREATE INDEX "user@balance_id_idx" ON shop."user" (balance_id);
CREATE INDEX "user@inventory_id_idx" ON shop."user" (inventory_id);

-- migrate:down
DROP SCHEMA IF EXISTS shop CASCADE;