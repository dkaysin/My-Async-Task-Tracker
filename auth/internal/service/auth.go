package service

import (
	"async_course/auth"
	"context"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

func (s *Service) Login(ctx context.Context, userID, passwordHash string) (string, error) {
	// verify that userID and password are correct
	ok, err := s.verifyAccountCredentials(ctx, userID, passwordHash)
	if err != nil {
		slog.Error("error while verifying credentials", "error", err)
		return "", err
	}
	if !ok {
		slog.Error("invalid credentials", "error", auth.ErrInvalidCredentials)
		return "", auth.ErrInvalidCredentials
	}

	// fetch claims (user's role) from db
	claims, err := s.getClaimsForAccount(ctx, userID)
	if err != nil {
		slog.Error("error while getting jwt claims", "error", err)
		return "", err
	}

	// prepare jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.signingKey))
	if err != nil {
		slog.Error("error while signing string", "error", err)
		return "", err
	}
	return tokenString, err
}

func (s *Service) verifyAccountCredentials(ctx context.Context, userID, passwordHash string) (bool, error) {
	var exists bool
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `SELECT COUNT(*) > 0 FROM accounts WHERE user_id = $1 AND password_hash = $2`
		return tx.QueryRow(context.Background(), q, userID, passwordHash).Scan(&exists)
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) getClaimsForAccount(ctx context.Context, userID string) (auth.JwtCustomClaims, error) {
	var role string
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `SELECT role FROM accounts WHERE user_id = $1`
		return tx.QueryRow(ctx, q, userID).Scan(&role)
	})
	return auth.JwtCustomClaims{
		UserID: userID,
		Role:   role,
	}, err
}
