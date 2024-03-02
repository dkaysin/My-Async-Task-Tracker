package service

import (
	schema "async_course/schema_registry"
	"async_course/task"

	"context"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskUserID struct {
	TaskID string  `json:"task_id"`
	UserID *string `json:"user_id"`
}

func (s *Service) GetTasksForAccount(ctx context.Context, userID string) ([]task.Task, error) {
	var tasks []task.Task
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT task_id, user_id, description, completed, created_at FROM tasks WHERE user_id = $1`
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

func (s *Service) CreateTask(ctx context.Context, description string) (string, *string, error) {
	taskID := uuid.New().String()

	userIDs, err := s.getActiveAccountsByRole(ctx, task.RoleDeveloper)
	if err != nil {
		return "", nil, err
	}
	var userID *string
	if !(len(userIDs) == 0) {
		userID = &userIDs[rand.IntN(len(userIDs))]
	}

	var createdTask task.Task
	err = s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO tasks (task_id, user_id, description, completed) VALUES ($1, $2, $3, False)
			RETURNING task_id, user_id, description, completed, created_at`
		rows, err := tx.Query(ctx, q, taskID, userID, description)
		if err != nil {
			return err
		}
		defer rows.Close()
		createdTask, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[task.Task])
		return err
	})
	if err != nil {
		return "", nil, err
	}
	if userID != nil {
		message := s.ew.SchemaRegistry.NewEventTaskAssigned(schema.Task(createdTask), nil)
		s.ew.TopicWriterTask.WriteMessage(message)
	}
	return taskID, userID, nil
}

func (s *Service) CompleteTask(ctx context.Context, taskID, userID string) error {
	var completedTask task.Task
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `UPDATE tasks SET completed = True, updated_at = NOW() WHERE task_id = $1 AND user_id = $2 AND completed = False
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
	if err != nil {
		return err
	}
	message := s.ew.SchemaRegistry.NewEventTaskCompleted(schema.Task(completedTask))
	s.ew.TopicWriterTask.WriteMessage(message)
	return nil
}

func (s *Service) AssignTasks(ctx context.Context) error {
	developerIDs, err := s.getActiveAccountsByRole(ctx, task.RoleDeveloper)
	if err != nil {
		return err
	}
	if len(developerIDs) == 0 {
		return task.ErrNoDevelopersAvailable
	}
	taskUserIDs, err := s.getNonCompletedTasks(ctx)
	if err != nil {
		return err
	}

	for _, taskUserID := range taskUserIDs {
		n := rand.IntN(len(developerIDs))
		developerID := developerIDs[n]
		// skip task assign if developer did not change
		if taskUserID.UserID != nil && *taskUserID.UserID == developerID {
			continue
		}

		var assignedTask task.Task
		err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
			var err error
			assignedTask, err = s.assignTaskTx(tx, ctx, taskUserID.TaskID, developerID)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		message := s.ew.SchemaRegistry.NewEventTaskAssigned(schema.Task(assignedTask), taskUserID.UserID)
		s.ew.TopicWriterTask.WriteMessage(message)
	}
	return nil
}

func (s *Service) getNonCompletedTasks(ctx context.Context) ([]TaskUserID, error) {
	var assignedTasks []TaskUserID
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT task_id, user_id FROM tasks WHERE NOT completed`
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()
		assignedTasks, err = pgx.CollectRows(rows, pgx.RowToStructByName[TaskUserID])
		return err
	})
	if err != nil {
		return nil, err
	}
	return assignedTasks, nil
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
