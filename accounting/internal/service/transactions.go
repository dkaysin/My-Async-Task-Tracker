package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// balance types
//
// accounts:
// - debit - user owes company money
// - credit - company owes user money
//
// cash:
// - debit - company's cash increases
// - credit - company's cash decreases
//
// profit:
// - debit - cost incurred
// - credit - revenue earned
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
type userAmount struct {
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Service) ProcessAssignTask(ctx context.Context, userID string) error {
	// process assign task transaction
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		return s.processTransactionTx(tx, ctx, newAssignTaskTransaction(userID))
	})
}

func (s *Service) ProcessCompleteTask(ctx context.Context, userID string) error {
	// process complete task transaction
	return s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		return s.processTransactionTx(tx, ctx, newCompleteTaskTransacion(userID))
	})
}

func (s *Service) ProcessCloseBalances(ctx context.Context) error {
	var processedUsers []userAmount
	err := s.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		// get list of users to whom we owe money
		usersWithDeficit, err := getUsersWithBalanceDeficitTx(tx, ctx, "accounts")
		if err != nil {
			return err
		}
		for _, userWithDeficit := range usersWithDeficit {
			userID := userWithDeficit.UserID
			deficit := userWithDeficit.Amount
			if deficit <= 0 {
				continue
			}
			slog.Info("processing user balance close", "user_id", userID, "deficit", deficit)

			// process balance close transaction
			err := s.processTransactionTx(tx, ctx, newBalanceCloseTransaction(userID, deficit))
			if err != nil {
				return err
			}

			// pay money to user and save payment timestamp
			processedAt, err := s.processPayment(userID, deficit)
			if err != nil {
				slog.Error("error while processing payment", "user_id", userID, "amount", deficit)
				return err
			}

			processedUsers = append(processedUsers, userAmount{userID, deficit, processedAt})
		}
		return nil
	})
	if err != nil {
		slog.Error("error while closing balances", "error", err)
		return err
	}
	// send event messages for payments made
	for _, user := range processedUsers {
		message := s.ew.SchemaRegistry.NewEventPaymentMade(user.UserID, user.Amount, user.Timestamp)
		err = s.ew.TopicWriterPayment.WriteMessage(message)
		if err != nil {
			slog.Error("error while sending payment message", "user_id", user.UserID, "processed_at", user.Timestamp, "error", err)
			return err
		}

	}
	return nil
}

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
		Source:            "balance_closed",
		Amount:            amount,
	}
}

type transactionFnTx func(pgx.Tx, context.Context, transaction) error

func (s *Service) processTransactionTx(tx pgx.Tx, ctx context.Context, t transaction) error {
	fns := []transactionFnTx{
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

func getUsersWithBalanceDeficitTx(tx pgx.Tx, ctx context.Context, balanceType string) ([]userAmount, error) {
	var userAmounts []userAmount
	q := `SELECT user_id, -(debit - credit) AS amount, NOW() as timestamp FROM balances WHERE balance_type = $1 AND debit - credit < 0`
	rows, err := tx.Query(ctx, q, balanceType)
	if err != nil {
		return userAmounts, err
	}
	defer rows.Close()
	userAmounts, err = pgx.CollectRows(rows, pgx.RowToStructByName[userAmount])
	return userAmounts, err
}
