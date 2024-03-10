--liquibase formatted sql

--changeset dmitry.k:3 endDelimiter:$$

CREATE TYPE balance_type AS ENUM ('accounts', 'cash', 'profit');

CREATE TABLE transactions(
    log_id TEXT NOT NULL PRIMARY KEY,
    transaction_id TEXT NOT NULL,
    balance_type balance_type NOT NULL,
    user_id TEXT,
    source TEXT NOT NULL,
    debit INT NOT NULL CHECK(debit >= 0),
    credit INT NOT NULL CHECK(credit >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE balances(
    balance_type balance_type NOT NULL,
    user_id TEXT,
    debit INT NOT NULL CHECK(debit >= 0),
    credit INT NOT NULL CHECK(credit >= 0),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (balance_type, user_id)
);

CREATE INDEX balances_balance_type_user_id_idx ON balances(balance_type, user_id);

CREATE TABLE accounting_accounts(
    user_id TEXT PRIMARY KEY,
    active BOOL NOT NULL,
    role TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE accounting_tasks(
    task_id TEXT PRIMARY KEY,
    user_id TEXT,
    description TEXT NOT NULL,
    completed BOOL NOT NULL DEFAULT False,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE accounting_dlq(
    message_headers TEXT,
    message_key TEXT,
    message_value BYTEA,
    error TEXT
);
