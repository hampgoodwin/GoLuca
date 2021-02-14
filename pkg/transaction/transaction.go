package transaction

import (
	"context"
	"errors"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
)

// Entry makes an entry to the accounting ledger
func (t *Transaction) Entry(ctx context.Context) error {
	if !t.balanced() {
		return errors.New("this transaction is not balanced")
	}

	transactionInsertStmt, err := data.DB.PrepareContext(ctx,
		`INSERT INTRO transaction(description) VALUES(?);`)
	result, err := transactionInsertStmt.ExecContext(ctx, t.Description)
	if err != nil {
		return err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, entry := range t.Entries {
		entryInsertStmt, err := data.DB.PrepareContext(ctx,
			`INSERT INTO entry(transaction_id, account_id, amount) VALUES(?, ?, ?)`)
		if err != nil {
			return err
		}

		entryInsertStmt.ExecContext(ctx, transactionID, entry.Account, entry.Amount)
	}

	return nil
}

// balanced checks that a transaction is balanced, that is to say that debits equal credits
func (t *Transaction) balanced() bool {
	sum := float64(0)
	for _, e := range t.Entries {
		sum += e.Amount
	}
	return sum == float64(0)
}
