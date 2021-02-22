package data

import (
	"context"

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
	if err := transactionSelectSmt.QueryRowContext(ctx, transactionID).Scan(
		&transaction.ID,
		&transaction.Description,
	); err != nil {
		return nil, errors.Wrap(err, "failed to scan row from transaction query result set")
	}
	return transaction, nil
}

// GetTransactions get's transactions paginaged by cursor and limit
func GetTransactions(ctx context.Context, cursor int64, limit int64) ([]transaction.Transaction, error) {
	transactionsSelectStmt, err := DB.PrepareContext(ctx,
		`SELECT id, description
FROM transaction
WHERE transaction.id > $1
LIMIT $2
;`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare transactions get statement")
	}
	rows, err := transactionsSelectStmt.QueryContext(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query db for transactions")
	}
	defer rows.Close()
	transactions := []transaction.Transaction{}
	for rows.Next() {
		transaction := transaction.Transaction{}
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan transactions results set")
		}
		transactions = append(transactions, transaction)
	}

	for i, transaction := range transactions {
		entries, err := GetEntriesByTransactionID(ctx, transaction.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to query db for entries when getting transactions")
		}
		transactions[i].Entries = append(transactions[i].Entries, entries...)
	}
	return transactions, nil
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
	if err := txInsertTransactionStmt.QueryRowContext(ctx, trans.Description).
		Scan(&transactionCreated.ID, &transactionCreated.Description); err != nil {
		return nil, errors.Wrap(err, "failed to scan transaction created return result set")
	}

	// insert the entries
	for _, entry := range trans.Entries {
		txInsertEntryStmt, err := tx.PrepareContext(ctx, `
INSERT INTO entry(transaction_id, account_id, amount) VALUES ($1, $2, $3)
RETURNING id, transaction_id, account_id, amount;`)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, errors.Wrap(err, "failed to rollback on failed transaction creation commit")
			}
			return nil, errors.Wrap(err, "failed to prepare entry insert statement")
		}
		entryCreated := transaction.Entry{}
		if err := txInsertEntryStmt.QueryRowContext(ctx, transactionCreated.ID, entry.AccountID, entry.Amount).
			Scan(&entryCreated.ID,
				&entryCreated.TransactionID,
				&entryCreated.AccountID,
				&entryCreated.Amount); err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction entry created return result set")
		}
		transactionCreated.Entries = append(transactionCreated.Entries, entryCreated)
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, errors.Wrap(err, "failed to rollback on failed transaction creation commit")
		}
		return nil, errors.Wrap(err, "failed to commit on transaction creation")
	}
	return transactionCreated, nil
}
