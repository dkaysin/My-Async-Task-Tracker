package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) ProcessAssignTask(ctx context.Context, userID string) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		return s.processTransactionTx(tx, ctx, newAssignTaskTransaction(userID))
	})
}

func (s *Service) ProcessCompleteTask(ctx context.Context, userID string) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		return s.processTransactionTx(tx, ctx, newCompleteTaskTransacion(userID))
	})
}

func (s *Service) ProcessBalanceClose(ctx context.Context, userID string) error {
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		balance, err := getOutstandingBalance(tx, ctx, balanceTypeAccounts, userID)
		if err != nil {
			return err
		}
		if balance < 0 {
			return s.processTransactionTx(tx, ctx, newBalanceCloseTransaction(userID, -balance))
		}
		return nil
	})
}

const (
	balanceTypeAccounts = "accounts"
	balanceTypeCash     = "cash"
	balanceTypeProfit   = "profit"
)

type transaction struct {
	ID                string
	BalanceTypeDebit  string
	BalanceTypeCredit string
	UserID            *string
	Source            string
	Amount            int
}

type transactinFnTx func(pgx.Tx, context.Context, transaction) error

func newAssignTaskTransaction(userID string) transaction {
	return transaction{
		ID:                uuid.New().String(),
		BalanceTypeDebit:  balanceTypeAccounts,
		BalanceTypeCredit: balanceTypeProfit,
		UserID:            &userID,
		Source:            "task_assigned",
		Amount:            priceAssignTask(),
	}
}

func newCompleteTaskTransacion(userID string) transaction {
	return transaction{
		ID:                uuid.New().String(),
		BalanceTypeDebit:  balanceTypeProfit,
		BalanceTypeCredit: balanceTypeAccounts,
		UserID:            &userID,
		Source:            "task_completed",
		Amount:            priceCompleteTask(),
	}
}

func newBalanceCloseTransaction(userID string, amount int) transaction {
	return transaction{
		ID:                uuid.New().String(),
		BalanceTypeDebit:  balanceTypeAccounts,
		BalanceTypeCredit: balanceTypeCash,
		UserID:            &userID,
		Source:            "balance_close",
		Amount:            amount,
	}
}

func (s *Service) processTransactionTx(tx pgx.Tx, ctx context.Context, t transaction) error {
	fns := []transactinFnTx{
		insertDebitTransactionTx,
		insertCreditTransactionTx,
		updateBalanceDebitTx,
		updateBalanceCreditTx,
	}
	for _, fn := range fns {
		err := fn(tx, ctx, t)
		if err != nil {
			slog.Error("error while processing transaction", "transaction", t)
			return err
		}
	}
	return nil
}

func insertDebitTransactionTx(tx pgx.Tx, ctx context.Context, t transaction) error {
	logID := uuid.New().String()
	q := `INSERT INTO transactions (log_id, transaction_id, balance_type, user_id, source, debit, credit, created_at) VALUES ($1, $2, $3, $4, $5, $6, 0, NOW())`
	_, err := tx.Exec(ctx, q, logID, t.ID, t.BalanceTypeDebit, t.UserID, t.Source, t.Amount)
	return err
}

func insertCreditTransactionTx(tx pgx.Tx, ctx context.Context, t transaction) error {
	logID := uuid.New().String()
	q := `INSERT INTO transactions (log_id, transaction_id, balance_type, user_id, source, debit, credit, created_at) VALUES ($1, $2, $3, $4, $5, 0, $6, NOW())`
	_, err := tx.Exec(ctx, q, logID, t.ID, t.BalanceTypeCredit, t.UserID, t.Source, t.Amount)
	return err
}

func updateBalanceDebitTx(tx pgx.Tx, ctx context.Context, t transaction) error {
	q := `INSERT INTO balances (balance_type, user_id, debit, credit, updated_at) VALUES ($1, $2, $3, 0, NOW())
        ON CONFLICT (balance_type, user_id) DO UPDATE SET debit = balances.debit + $3`
	_, err := tx.Exec(ctx, q, t.BalanceTypeDebit, t.UserID, t.Amount)
	return err
}

func updateBalanceCreditTx(tx pgx.Tx, ctx context.Context, t transaction) error {
	q := `INSERT INTO balances (balance_type, user_id, debit, credit, updated_at) VALUES ($1, $2, 0, $3, NOW())
        ON CONFLICT (balance_type, user_id) DO UPDATE SET credit = balances.credit + $3`
	_, err := tx.Exec(ctx, q, t.BalanceTypeCredit, t.UserID, t.Amount)
	return err
}

func getOutstandingBalance(tx pgx.Tx, ctx context.Context, balanceType, userID string) (int, error) {
	q := `SELECT debit - credit FROM balances WHERE balance_type = $1, user_id = $2`
	var balance int
	err := tx.QueryRow(ctx, q, balanceType, userID).Scan(&balance)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return balance, err
}
