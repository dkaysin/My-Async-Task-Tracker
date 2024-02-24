package service

import (
	"async_course/task"
	"context"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) GetTasks(ctx context.Context, userID string) ([]task.Task, error) {
	var tasks []task.Task
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT task_id, description, completed, created_at WHERE user_id = $1`
		rows, err := tx.Query(ctx, q, userID)
		if err != nil {
			return err
		}
		defer rows.Close()
		tasks, err = pgx.CollectRows(rows, pgx.RowToStructByName[task.Task])
		return err
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) CreateTask(ctx context.Context, description string) (string, error) {
	taskID := uuid.New().String()
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO tasks (task_id, description, completed) VALUES ($1, $2, False)`
		_, err := tx.Exec(ctx, q, taskID, description)
		return err
	})
	if err != nil {
		return "", err
	}
	return taskID, nil
}

func (s *Service) CompleteTask(ctx context.Context, taskID, userID string) error {
	var completedTask task.Task
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `UPDATE tasks SET completed = True, updated_at = NOW() WHERE task_id = $1 AND user_id = $2
			RETURNING task_id, user_id, description, completed, created_at`
		rows, err := tx.Query(ctx, q, taskID, userID)
		if err != nil {
			return err
		}
		defer rows.Close()
		completedTask, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[task.Task])
		if err == pgx.ErrNoRows {
			return task.ErrTaskNotFound
		}
		return err
	})
	event := task.NewEventTaskCompleted(completedTask)
	s.ew.TopicWriterTask.WriteJSON(context.Background(), event.Key, event.Value)
	return err
}

func (s *Service) AssignTasks(ctx context.Context) error {
	taskIDs, err := s.getUnassignedTasks(ctx)
	if err != nil {
		return err
	}
	developerIDs, err := s.getActiveAccountsByRole(ctx, task.RoleDeveloper)
	if err != nil {
		return err
	}

	for _, taskID := range taskIDs {
		n := rand.IntN(len(developerIDs))
		developerID := developerIDs[n]
		var assignedTask task.Task
		err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
			var err error
			assignedTask, err = s.assignTaskTx(tx, ctx, taskID, developerID)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		event := task.NewEventTaskCompleted(assignedTask)
		s.ew.TopicWriterTask.WriteJSON(context.Background(), event.Key, event.Value)
	}
	return nil
}

func (s *Service) getUnassignedTasks(ctx context.Context) ([]string, error) {
	var taskIDs []string
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT task_id FROM tasks WHERE user_id IS NULL AND NOT completed`
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()
		taskIDs, err = pgx.CollectRows(rows, pgx.RowTo[string])
		return err
	})
	if err != nil {
		return nil, err
	}
	return taskIDs, nil
}

func (s *Service) assignTaskTx(tx pgx.Tx, ctx context.Context, taskID, userID string) (task.Task, error) {
	q := `UPDATE tasks SET user_id = $2, updated_at = NOW() WHERE task_id = $1
		RETURNING task_id, user_id, description, completed, created_at`
	rows, err := tx.Query(ctx, q, taskID, userID)
	if err != nil {
		return task.Task{}, err
	}
	defer rows.Close()
	assignedTask, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[task.Task])
	if err == pgx.ErrNoRows {
		return assignedTask, pgx.ErrNoRows
	}
	return assignedTask, err
}
