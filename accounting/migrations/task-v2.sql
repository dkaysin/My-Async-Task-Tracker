--liquibase formatted sql

--changeset dmitry.k:6 endDelimiter:$$

ALTER TABLE accounting_tasks ADD COLUMN jira_id TEXT;

