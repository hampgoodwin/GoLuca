package data

import (
	"context"
	"database/sql"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/jackc/pgx/v4"
)

// GetTransaction get's a transaction record, without it's entries, by the transaction ID
func GetTransaction(ctx context.Context, transactionID int64) (*transaction.Transaction, error) {
	transaction := &transaction.Transaction{}
	if err := DBPool.QueryRow(ctx, `SELECT id, description
FROM transaction
WHERE id=$1
;`, transactionID).Scan(
		&transaction.ID,
		&transaction.Description,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WrapFlag(err, "scanning transaction result row", errors.NotFound)
		}
		return nil, errors.Wrap(err, "scanning row from transaction query result set")
	}
	return transaction, nil
}

// GetTransactions get's transactions paginaged by cursor and limit
func GetTransactions(ctx context.Context, cursor int64, limit int64) ([]transaction.Transaction, error) {
	rows, err := DBPool.Query(ctx, `SELECT id, description
FROM transaction
WHERE transaction.id > $1
LIMIT $2
;`, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying database for transactions")
	}
	defer rows.Close()
	transactions := []transaction.Transaction{}
	for rows.Next() {
		transaction := transaction.Transaction{}
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
		); err != nil {
			return nil, errors.Wrap(err, "scanning transactions results set")
		}
		transactions = append(transactions, transaction)
	}

	for i, transaction := range transactions {
		entries, err := GetEntriesByTransactionID(ctx, transaction.ID)
		if err != nil {
			return nil, errors.Wrap(err, "querying databse for entries while getting transactions")
		}
		transactions[i].Entries = append(transactions[i].Entries, entries...)
	}
	return transactions, nil
}

// CreateTransaction creates a transaction and associated entries in a single transaction
func CreateTransaction(ctx context.Context, trans *transaction.Transaction) (*transaction.Transaction, error) {
	// get a db-transaction
	tx, err := DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "starting database transaction for creating transaction")
	}
	transactionCreated := &transaction.Transaction{}
	if err := tx.QueryRow(ctx, `
INSERT INTO transaction (description) VALUES ($1)
RETURNING id, description;`, trans.Description).
		Scan(&transactionCreated.ID, &transactionCreated.Description); err != nil {
		return nil, errors.Wrap(err, "scanning transaction created return result set")
	}

	// insert the entries
	for _, entry := range trans.Entries {
		entryCreated := transaction.Entry{}
		if err := tx.QueryRow(ctx, `
INSERT INTO entry(transaction_id, account_id, amount) VALUES ($1, $2, $3)
RETURNING id, transaction_id, account_id, amount;`,
			transactionCreated.ID, entry.AccountID, entry.Amount).
			Scan(&entryCreated.ID,
				&entryCreated.TransactionID,
				&entryCreated.AccountID,
				&entryCreated.Amount); err != nil {
			return nil, errors.Wrap(err, "scanning transaction entry created return result set")
		}
		transactionCreated.Entries = append(transactionCreated.Entries, entryCreated)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "committing on transaction creation")
	}
	return transactionCreated, nil
}
