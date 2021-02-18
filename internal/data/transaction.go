package data

import (
	"context"
	"fmt"

	"github.com/abelgoodwin1988/GoLuca/pkg/transaction"
	"github.com/pkg/errors"
)

// GetTransaction get's a transaction record, without it's entries, by the transaction ID
func GetTransaction(ctx context.Context, transactionID int64) (*transaction.Transaction, error) {
	transactionSelectSmt, err := DB.PrepareContext(ctx,
		`SELECT id, description
FROM transaction
WHERE id=$1
;`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare account get statement")
	}
	transaction := &transaction.Transaction{}
	transactionSelectSmt.QueryRowContext(ctx, transactionID).Scan(
		&transaction.ID,
		&transaction.Description,
	)
	return transaction, nil
}

// CreateTransaction creates a transaction and associated entries in a single transaction
func CreateTransaction(ctx context.Context, trans *transaction.Transaction) (*transaction.Transaction, error) {
	// get a db-transaction
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start db-transaction for creating transaction")
	}

	// insert the transaction
	txInsertTransactionStmt, err := tx.PrepareContext(ctx, `
INSERT INTO transaction (description) VALUES ($1)
RETURNING id, description;`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare transaction insert statement")
	}
	transactionCreated := &transaction.Transaction{}
	txInsertTransactionStmt.QueryRowContext(ctx, trans.Description).
		Scan(&transactionCreated.ID, transactionCreated.Description)

	// insert the entries
	for _, entry := range trans.Entries {
		txInsertEntryStmt, err := tx.PrepareContext(ctx, `
INSERT INTO entry(transaction_id, account_id, amount) VALUES ($1, $2, $3)
RETURNING id, transaction_id, account_id, amount;`)
		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare entry insert statement")
		}
		entryCreated := transaction.Entry{}
		txInsertEntryStmt.QueryRowContext(ctx, transactionCreated.ID, entry.AccountID, entry.Amount).
			Scan(&entryCreated.ID, &entryCreated.TransactionID, &entryCreated.AccountID, &entryCreated.Amount)
		fmt.Println(entryCreated)
		transactionCreated.Entries = append(transactionCreated.Entries, entryCreated)
	}
	tx.Commit()
	return transactionCreated, nil
}
