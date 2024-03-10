package service

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Service) getActiveAccountsByRole(ctx context.Context, role string) ([]string, error) {
	var userIDs []string
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT user_id FROM task_accounts WHERE role = $1 AND active`
		rows, err := tx.Query(ctx, q, role)
		if err != nil {
			return err
		}
		defer rows.Close()
		userIDs, err = pgx.CollectRows(rows, pgx.RowTo[string])
		return err
	})
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (s *Service) UpsertAccountRole(ctx context.Context, userID string, active bool, role string) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO task_accounts (user_id, active, role) VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET active = $2, role = $3, updated_at = NOW()`
		_, err := tx.Exec(ctx, q, userID, active, role)
		return err
	})
}
