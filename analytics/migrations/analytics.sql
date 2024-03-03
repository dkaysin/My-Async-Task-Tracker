--liquibase formatted sql

--changeset dmitry.k:4 endDelimiter:$$

CREATE TABLE analytics(
    source TEXT NOT NULL,
    user_id TEXT,
    revenue INT NOT NULL,
    cost INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

