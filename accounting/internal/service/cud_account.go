package service

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Service) UpsertAccountRole(ctx context.Context, userID string, active bool, role string) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO accounting_accounts (user_id, active, role) VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET active = $2, role = $3, updated_at = NOW()`
		_, err := tx.Exec(ctx, q, userID, active, role)
		return err
	})
}
