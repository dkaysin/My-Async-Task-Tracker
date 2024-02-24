package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pgPool *pgxpool.Pool
}

func NewDatabase(pgConnString string) (*Database, error) {
	pgPool, err := pgxpool.New(context.Background(), pgConnString)
	if err != nil {
		slog.Error("failed to create pg pool", "error", err)
		return nil, err
	}
	err = pgPool.Ping(context.Background())
	if err != nil {
		slog.Error("failed to ping database", "error", err)
		return nil, err
	}
	return &Database{pgPool}, nil
}

type txFunc func(tx pgx.Tx) error

func (db *Database) ExecuteTx(ctx context.Context, fn txFunc) error {
	tx, err := db.pgPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = fn(tx)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error executing transaction", "error", err)
		return err
	}

	return tx.Commit(ctx)
}

func (db *Database) Close() {
	db.pgPool.Close()
}
