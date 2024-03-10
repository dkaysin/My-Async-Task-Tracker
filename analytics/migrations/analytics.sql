--liquibase formatted sql

--changeset dmitry.k:4 endDelimiter:$$

CREATE TABLE analytics(
    source TEXT NOT NULL,
    user_id TEXT,
    revenue INT NOT NULL,
    cost INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE analytics_accounts(
    user_id TEXT PRIMARY KEY,
    active BOOL NOT NULL,
    role TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE analytics_tasks(
    task_id TEXT PRIMARY KEY,
    user_id TEXT,
    description TEXT NOT NULL,
    completed BOOL NOT NULL DEFAULT False,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
