package data

import (
	"context"

	"github.com/abelgoodwin1988/GoLuca/pkg/transaction"
	"github.com/pkg/errors"
)

// GetEntries gets a paginated result of db entries
func GetEntries(ctx context.Context, cursor int64, limit int64) ([]transaction.Entry, error) {
	getEntriesStmt, err := DB.Prepare(`SELECT id, transaction_id, account_id, amount FROM entry WHERE id > $1 LIMIT $2;`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entries on db query")
	}
	rows, err := getEntriesStmt.QueryContext(ctx, cursor, limit)
	var entries []transaction.Entry
	for rows.Next() {
		entry := transaction.Entry{}
		rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.AccountID,
			&entry.Amount,
		)
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetEntriesByTransactionID gets entries by transaction ID
func GetEntriesByTransactionID(ctx context.Context, transactionID int64) ([]transaction.Entry, error) {
	getEntriesStmt, err := DB.Prepare(`SELECT id, transaction_id, account_id, amount FROM entry WHERE transaction_id=$1`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entries on db query")
	}
	rows, err := getEntriesStmt.QueryContext(ctx, transactionID)
	var entries []transaction.Entry
	for rows.Next() {
		entry := transaction.Entry{}
		rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.AccountID,
			&entry.Amount,
		)
		entries = append(entries, entry)
	}
	return entries, nil
}
