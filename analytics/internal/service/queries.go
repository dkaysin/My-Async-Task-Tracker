package service

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DevelopersReport struct {
	CountDeficit int `json:"count_deficit"`
}

type ProfitReportEntry struct {
	Date    string `json:"date"`
	Revenue int    `json:"revenue"`
	Cost    int    `json:"cost"`
	Profit  int    `json:"profit"`
}

type RevenueSourceReportEntry struct {
	Revenue int    `json:"revenue"`
	Source  string `json:"source"`
	UserID  string `json:"user_id"`
}

type RevenueSourceReportItem struct {
	Date string `json:"date"`
	RevenueSourceReportEntry
}

func (s *Service) GetDevelopersReport(ctx context.Context) (DevelopersReport, error) {
	var res DevelopersReport
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT count(DISTINCT user_id) AS count_deficit FROM analytics WHERE revenue - cost > 0`
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()
		res, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[DevelopersReport])
		return err
	})
	return res, err
}

func (s *Service) GetProfitReport(ctx context.Context) ([]ProfitReportEntry, error) {
	var entries []ProfitReportEntry
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS date, SUM(revenue) AS revenue, SUM(cost) AS cost, SUM(revenue - cost) AS profit
            FROM analytics
            GROUP BY date
            ORDER BY date DESC`
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()
		entries, err = pgx.CollectRows(rows, pgx.RowToStructByName[ProfitReportEntry])
		return err
	})
	return entries, err
}

func (s *Service) GetRevenueSourceReport(ctx context.Context) (map[string][]RevenueSourceReportEntry, error) {
	var items []RevenueSourceReportItem
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `WITH t AS(
                SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS date, revenue, source, user_id,
                MAX(revenue) OVER (PARTITION BY TO_CHAR(created_at, 'YYYY-MM-DD')) AS max_revenue
                FROM analytics)
            SELECT date, revenue, source, user_id FROM t
            WHERE revenue = max_revenue`
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()
		items, err = pgx.CollectRows(rows, pgx.RowToStructByName[RevenueSourceReportItem])
		return err
	})
	entries := map[string][]RevenueSourceReportEntry{}
	for _, item := range items {
		date := item.Date
		entry := item.RevenueSourceReportEntry
		itemsForDate, ok := entries[date]
		if !ok {
			itemsForDate = []RevenueSourceReportEntry{}
		}
		entries[date] = append(itemsForDate, entry)
	}
	return entries, err
}
