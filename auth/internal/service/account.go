package service

import (
	"async_course/auth"

	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) CreateAccount(ctx context.Context, name, passwordHash, role string) (string, error) {
	userID := uuid.New().String()
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO accounts (user_id, name, password_hash, role) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(ctx, q, userID, name, passwordHash, role)
		return err
	})
	if err != nil {
		return "", err
	}
	message := s.ew.SchemaRegistry.NewEventAccountCreated(userID, role)
	s.ew.TopicWriterAccount.WriteMessage(message)
	return userID, nil
}

func (s *Service) ChangeAccountRole(ctx context.Context, userID, newRole string) error {
	var active bool
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `UPDATE accounts SET role = $2, updated_at = NOW() WHERE user_id = $1
			RETURNING active`
		err := tx.QueryRow(ctx, q, userID, newRole).Scan(&active)
		if err == pgx.ErrNoRows {
			return auth.ErrAccountNotFound
		}
		return err
	})
	if err != nil {
		return err
	}
	message := s.ew.SchemaRegistry.NewEventAccountUpdated(userID, newRole, active)
	s.ew.TopicWriterAccount.WriteMessage(message)
	return nil
}
