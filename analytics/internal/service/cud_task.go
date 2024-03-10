package service

import (
	schema "async_course/schema_registry"

	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Service) UpsertTask(ctx context.Context, task schema.Task) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO analytics_tasks (task_id, user_id, description, completed) VALUES ($1, $2, $3, $4)
			ON CONFLICT (task_id) DO UPDATE SET user_id = $2, description = $3, completed = $4, updated_at = NOW()`
		_, err := tx.Exec(ctx, q, task.TaskID, task.UserID, task.Description, task.Completed)
		return err
	})
}
