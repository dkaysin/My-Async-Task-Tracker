--liquibase formatted sql

--changeset dmitry.k:7 endDelimiter:$$

ALTER TABLE analytics_tasks ADD COLUMN jira_id TEXT;
