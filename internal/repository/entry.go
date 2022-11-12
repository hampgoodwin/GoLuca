package repository

import (
	"context"
	"fmt"

	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
	"github.com/jackc/pgx/v4"
)

// GetEntriesByTransactionID gets entries by transaction ID
func getEntriesByTransactionID(ctx context.Context, tx pgx.Tx, transactionID string) ([]Entry, error) {
	rows, err := tx.Query(ctx,
		`SELECT id, transaction_id, description, debit_account, credit_account, amount_value, amount_currency, created_at
		FROM entry
		WHERE transaction_id=$1
		;`,
		transactionID)
	if err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotKnown,
			fmt.Sprintf("fetching entries by transaction id %q from datastore", transactionID))
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		entry := Entry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.Description,
			&entry.DebitAccount,
			&entry.CreditAccount,
			&entry.Amount.Value,
			&entry.Amount.Currency,
			&entry.CreatedAt,
		); err != nil {
			return nil, errors.WithErrorMessage(err, errors.NotKnown, "scanning entry result row")
		}
		entries = append(entries, entry)
	}

	if err := validate.Validate(entries); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating entries fetched from database")
	}

	return entries, nil
}
