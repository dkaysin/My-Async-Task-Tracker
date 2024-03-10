package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type AnalyticsEntry struct {
	UserID    *string
	Source    string
	Revenue   int
	Cost      int
	CreatedAt time.Time
}

func (s *Service) ProcessTaskAssigned(ctx context.Context, userID *string, source string, revenue int, createdAt time.Time) error {
	return s.insertAnalyticsEntry(ctx, AnalyticsEntry{
		UserID:    userID,
		Source:    source,
		Revenue:   revenue,
		Cost:      0,
		CreatedAt: createdAt,
	})
}

func (s *Service) ProcessTaskCompleted(ctx context.Context, userID *string, source string, cost int, createdAt time.Time) error {
	return s.insertAnalyticsEntry(ctx, AnalyticsEntry{
		UserID:    userID,
		Source:    source,
		Revenue:   0,
		Cost:      cost,
		CreatedAt: createdAt,
	})
}

func (s *Service) insertAnalyticsEntry(ctx context.Context, e AnalyticsEntry) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO analytics (source, user_id, revenue, cost, created_at) VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.Exec(ctx, q, e.Source, e.UserID, e.Revenue, e.Cost, e.CreatedAt)
		return err
	})
}
