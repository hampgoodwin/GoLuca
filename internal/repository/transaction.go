package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetTransaction get's a transaction record, without it's entries, by the transaction ID
func (r *Repository) GetTransaction(ctx context.Context, transactionID string) (Transaction, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "repository.GetTransaction", trace.WithAttributes(
		attribute.String("transaction_id", transactionID),
	))
	defer span.End()

	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		errors.Join(ierrors.ErrNotKnown, fmt.Errorf("beginning get transaction db transaction: %w", err))
	}
	returning := Transaction{}
	if err = tx.QueryRow(ctx,
		`SELECT id, description, created_at
		FROM transaction
		WHERE id=$1
		;`, transactionID).Scan(
		&returning.ID,
		&returning.Description,
		&returning.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return returning, errors.Join(fmt.Errorf("transaction %q not found: %w", transactionID, err), ierrors.ErrNotFound)
		}
		return returning, errors.Join(fmt.Errorf("scanning transaction result row: %w", err), ierrors.ErrNotKnown)
	}

	returning.Entries, err = getEntriesByTransactionID(ctx, tx, transactionID)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return returning, fmt.Errorf("rolling back on fetching transaction entries error: %w", err)
		}
		return returning, fmt.Errorf("getting entries by transaction %q: %w", transactionID, err)
	}

	if err := validate.Validate(returning); err != nil {
		return returning, errors.Join(fmt.Errorf("validating transaction fetched from database: %w", err), ierrors.ErrNotValidInternalData)
	}

	if err := tx.Commit(ctx); err != nil {
		return returning, errors.Join(fmt.Errorf("committing get transaction query: %w"), ierrors.ErrNotKnown)
	}

	return returning, nil
}

// ListTransactions get's transactions paginated by cursor and limit
func (r *Repository) ListTransactions(ctx context.Context, transactionID string, createdAt time.Time, limit uint64) ([]Transaction, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "repository.ListTransaction", trace.WithAttributes(
		attribute.String("cursor.id", transactionID),
		attribute.String("cursor.created_at", createdAt.String()),
		attribute.Int64("limit", int64(limit)),
	))
	defer span.End()

	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.Join(fmt.Errorf("beginning get transactions db transactoion: %w", err), ierrors.ErrNotKnown)
	}

	query := `SELECT id, description, created_at
		FROM transaction
		WHERE 1=1`
	var params []interface{}
	if transactionID != "" && !createdAt.IsZero() {
		params = append(params, transactionID)
		query += fmt.Sprintf(" AND transaction.id <= $%d", len(params))
		params = append(params, createdAt)
		query += fmt.Sprintf(" AND created_at <= $%d", len(params))
	}
	query += " ORDER BY created_at DESC"
	if limit != 0 {
		params = append(params, limit)
		query += fmt.Sprintf(" LIMIT $%d", len(params))
	}
	query += ";"

	rows, err := tx.Query(ctx, query, params...)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("fetching transactions from database: %w", err), ierrors.ErrNotKnown)
	}
	defer rows.Close()
	returning := []Transaction{}
	for rows.Next() {
		transaction := Transaction{}
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.CreatedAt,
		); err != nil {
			return nil, errors.Join(fmt.Errorf("scanning transaction result row: %w", err), ierrors.ErrNotKnown)
		}
		returning = append(returning, transaction)
	}

	for i, transaction := range returning {
		entries, err := getEntriesByTransactionID(ctx, tx, transaction.ID)
		if err != nil {
			return nil, fmt.Errorf("getting entries by transaction id: %w", err)
		}
		returning[i].Entries = append(returning[i].Entries, entries...)
	}

	if err := validate.Validate(returning); err != nil {
		return nil, errors.Join(fmt.Errorf("validating transactions fetched from database: %w", err), ierrors.ErrNotValidInternalData)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Join(fmt.Errorf("committing get transactions transaction: %w", err), ierrors.ErrNotKnown)
	}

	return returning, nil
}

// CreateTransaction creates a transaction and associated entries in a single transaction
func (r *Repository) CreateTransaction(ctx context.Context, create Transaction) (Transaction, error) {
	// get a db-transaction
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Transaction{}, errors.Join(fmt.Errorf("beginning create transaction db transaction: %w", err), ierrors.ErrNotKnown)
	}
	returning := Transaction{}
	if err := tx.QueryRow(ctx,
		`INSERT INTO transaction (id, description, created_at) VALUES ($1, $2, $3)
		RETURNING id, description, created_at
		;`, create.ID, create.Description, create.CreatedAt).Scan(
		&returning.ID, &returning.Description, &returning.CreatedAt,
	); err != nil {
		return returning, errors.Join(fmt.Errorf("scanning transaction returned from insert: %w", err), ierrors.ErrNotKnown)
	}

	// insert the entries
	for _, entry := range create.Entries {
		returningEntry := Entry{}
		if err := tx.QueryRow(ctx,
			`INSERT INTO entry(id, transaction_id, description, debit_account, credit_account, amount_value, amount_currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, transaction_id, description, debit_account, credit_account, amount_value, amount_currency, created_at;`,
			entry.ID, returning.ID, entry.Description, entry.DebitAccount, entry.CreditAccount, entry.Amount.Value, entry.Amount.Currency, entry.CreatedAt).Scan(
			&returningEntry.ID,
			&returningEntry.TransactionID,
			&returningEntry.Description,
			&returningEntry.DebitAccount,
			&returningEntry.CreditAccount,
			&returningEntry.Amount.Value,
			&returningEntry.Amount.Currency,
			&returningEntry.CreatedAt,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.ForeignKeyViolation {
					return returning, errors.Join(fmt.Errorf("inserting entry with key constraint on one of transaction id, debit account, or credit account: %w", err), ierrors.ErrNoRelationshipFound)
				}
			}
			return returning, errors.Join(fmt.Errorf("scanning entry result row: %w", err), ierrors.ErrNotKnown)
		}
		returning.Entries = append(returning.Entries, returningEntry)
	}

	if err := validate.Validate(returning); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return returning, errors.Join(fmt.Errorf("rolling back create transaction on validating transaction returned from create: %w", err), ierrors.ErrNotKnown)
		}
		return returning, errors.Join(fmt.Errorf("validating transaction returned from creation: %w", err), ierrors.ErrNotValidInternalData)
	}

	if err := tx.Commit(ctx); err != nil {
		return returning, errors.Join(fmt.Errorf("committing create transaction transaction: %w", err), ierrors.ErrNotKnown)
	}
	return returning, nil
}
