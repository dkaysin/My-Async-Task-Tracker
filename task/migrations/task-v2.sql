--liquibase formatted sql

--changeset dmitry.k:5 endDelimiter:$$

ALTER TABLE tasks ADD COLUMN jira_id TEXT;

