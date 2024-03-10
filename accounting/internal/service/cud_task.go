package service

import (
	v1 "async_course/schema_registry/schemas/v1"
	v2 "async_course/schema_registry/schemas/v2"

	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Service) UpsertTaskV1(ctx context.Context, task v1.Task) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO accounting_tasks (task_id, user_id, description, completed) VALUES ($1, $2, $3, $4)
			ON CONFLICT (task_id) DO UPDATE SET user_id = $2, description = $3, completed = $4, updated_at = NOW()`
		_, err := tx.Exec(ctx, q, task.TaskID, task.UserID, task.Description, task.Completed)
		return err
	})
}

func (s *Service) UpsertTaskV2(ctx context.Context, task v2.Task) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO accounting_tasks (task_id, user_id, description, completed, jira_id) VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (task_id) DO UPDATE SET user_id = $2, description = $3, completed = $4, jira_id = $5, updated_at = NOW()`
		_, err := tx.Exec(ctx, q, task.TaskID, task.UserID, task.Description, task.Completed, task.JiraID)
		return err
	})
}
