package repository

import (
	"context"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

// GetEntries gets a paginated result of db entries
func (r *Repository) GetEntries(ctx context.Context, id string, createdAt time.Time, limit uint64) ([]transaction.Entry, error) {
	rows, err := r.Database.Query(ctx,
		`SELECT id, transaction_id, debit_account, credit_account, amount_value, amount_currency, created_at
		FROM entry
		WHERE id <= $1
			AND created_at <= $2
		ORDER BY created_at DESC
		LIMIT $2
		;`,
		id, createdAt, limit)
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
			&entry.DebitAccount,
			&entry.CreditAccount,
			&entry.Amount.Value,
			&entry.Amount.Currency,
			&entry.CreatedAt,
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
		`SELECT id, transaction_id, debit_account, credit_account, amount_value, amount_currency, created_at
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
			&entry.DebitAccount,
			&entry.CreditAccount,
			&entry.Amount.Value,
			&entry.Amount.Currency,
			&entry.CreatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning row from entries by transaction id query results set")
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
