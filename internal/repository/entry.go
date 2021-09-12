package repository

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

// GetEntries gets a paginated result of db entries
func (r *Repository) GetEntries(ctx context.Context, cursor string, limit uint64) ([]transaction.Entry, error) {
	rows, err := r.Database.Query(ctx,
		`SELECT id, transaction_id, account_id, amount
		FROM entry
		WHERE id > $1
		ORDER BY created_at
		LIMIT $2
		;`,
		cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting entries from database")
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
			return nil, errors.Wrap(err, "scanning row from entries query results set")
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetEntriesByTransactionID gets entries by transaction ID
func (r *Repository) GetEntriesByTransactionID(ctx context.Context, transactionID string) ([]transaction.Entry, error) {
	rows, err := r.Database.Query(ctx,
		`SELECT id, transaction_id, account_id, amount, created_at
		FROM entry
		WHERE transaction_id=$1
		;`,
		transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting entries by transaction id from database")
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
			&entry.CreatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning row from entries by transaction id query results set")
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
