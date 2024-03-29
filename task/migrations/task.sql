--liquibase formatted sql

--changeset dmitry.k:2 endDelimiter:$$

CREATE TABLE task_accounts(
    user_id TEXT PRIMARY KEY,
    active BOOL NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE tasks(
    task_id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES task_accounts(user_id) DEFAULT NULL,
    description TEXT NOT NULL,
    completed BOOL NOT NULL DEFAULT False,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);