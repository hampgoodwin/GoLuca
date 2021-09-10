package data

import (
	"context"

	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/pkg/errors"
)

// GetEntries gets a paginated result of db entries
func GetEntries(ctx context.Context, cursor int64, limit int64) ([]transaction.Entry, error) {
	rows, err := DBPool.Query(ctx, `SELECT id, transaction_id, account_id, amount FROM entry WHERE id > $1 LIMIT $2;`, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entries from db")
	}
	defer rows.Close()
	var entries []transaction.Entry
	for rows.Next() {
		entry := transaction.Entry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.AccountID,
			&entry.Amount,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan row from entries query results set")
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetEntriesByTransactionID gets entries by transaction ID
func GetEntriesByTransactionID(ctx context.Context, transactionID int64) ([]transaction.Entry, error) {
	rows, err := DBPool.Query(ctx, `SELECT id, transaction_id, account_id, amount FROM entry WHERE transaction_id=$1`, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "failured getting entries from db")
	}
	defer rows.Close()
	var entries []transaction.Entry
	for rows.Next() {
		entry := transaction.Entry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.AccountID,
			&entry.Amount,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan row from entries by transaction id query results set")
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
