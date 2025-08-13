package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetEntriesByTransactionID gets entries by transaction ID
func getEntriesByTransactionID(ctx context.Context, tx pgx.Tx, transactionID string) ([]Entry, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "repository.ListTransaction.getEntriesByTransactionID", trace.WithAttributes(
		attribute.String("transaction_id", transactionID),
		attribute.Int64("db.PID", int64(tx.Conn().PgConn().PID())),
	))
	defer span.End()

	rows, err := tx.Query(ctx,
		`SELECT id, transaction_id, description, debit_account, credit_account, amount_value, amount_currency, created_at
		FROM entry
		WHERE transaction_id=$1
		;`,
		transactionID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("fetching entries by transaction id %q from datastore: %w", transactionID, err), ierrors.ErrNotKnown)
	}
	defer rows.Close()
	var entries []Entry
	entryIDs := []string{}
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
			return nil, errors.Join(fmt.Errorf("scanning entry result row: %w", err), ierrors.ErrNotKnown)
		}
		entries = append(entries, entry)
		entryIDs = append(entryIDs, entry.ID)
	}
	span.SetAttributes(attribute.StringSlice("entry_ids", entryIDs))

	if err := validate.Validate(entries); err != nil {
		return nil, errors.Join(fmt.Errorf("validating entries fetched from database: %w", err), ierrors.ErrNotValidInternalData)
	}

	return entries, nil
}
