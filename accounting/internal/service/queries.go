package service

import (
	"async_course/accounting"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type BalanceSummary struct {
	Balance int `json:"balance"`
}

type BalanceLogEntry struct {
	TransactionID string    `json:"transaction_id"`
	Amount        int       `json:"amount"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProfitLogEntry struct {
	Date    string `json:"date"`
	Revenue int    `json:"revenue"`
	Cost    int    `json:"cost"`
	Profit  int    `json:"profit"`
}

func (s *Service) GetBalanceSummary(ctx context.Context, userID string) (BalanceSummary, error) {
	var balanceSummary BalanceSummary
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT credit - debit AS balance FROM balances WHERE user_id = $1 AND balance_type = 'accounts'`
		rows, err := tx.Query(ctx, q, userID)
		if err != nil {
			return err
		}
		defer rows.Close()
		balanceSummary, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[BalanceSummary])
		if err == pgx.ErrNoRows {
			return accounting.ErrUnknownUser
		}
		return err
	})
	return balanceSummary, err
}

func (s *Service) GetBalanceLog(ctx context.Context, userID string) ([]BalanceLogEntry, error) {
	var balanceLog []BalanceLogEntry
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT transaction_id, credit - debit AS amount, source, created_at FROM transactions
            WHERE user_id = $1 AND balance_type = 'accounts' ORDER BY created_at DESC`
		rows, err := tx.Query(ctx, q, userID)
		if err != nil {
			return err
		}
		defer rows.Close()
		balanceLog, err = pgx.CollectRows(rows, pgx.RowToStructByName[BalanceLogEntry])
		return err
	})
	return balanceLog, err
}

func (s *Service) GetProfitLog(ctx context.Context) ([]ProfitLogEntry, error) {
	var profitLog []ProfitLogEntry
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		var err error
		q := `SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS date, SUM(credit) AS revenue, SUM(debit) AS cost, SUM(credit - debit) AS profit
            FROM transactions WHERE balance_type = 'profit'
            GROUP BY date
            ORDER BY date DESC`
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()
		profitLog, err = pgx.CollectRows(rows, pgx.RowToStructByName[ProfitLogEntry])
		return err
	})
	return profitLog, err
}
