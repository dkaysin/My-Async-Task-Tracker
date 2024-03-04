--liquibase formatted sql

--changeset dmitry.k:1 endDelimiter:$$

CREATE TYPE account_role AS ENUM ('developer', 'admin', 'manager', 'accountant');

CREATE TABLE accounts(
    user_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    active BOOL NOT NULL DEFAULT True,
    
    role account_role NOT NULL,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO accounts (user_id, name, role, password_hash) VALUES ('admin', 'admin', 'admin', '12345');
