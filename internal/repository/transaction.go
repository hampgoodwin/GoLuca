package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

// GetTransaction get's a transaction record, without it's entries, by the transaction ID
func (r *Repository) GetTransaction(ctx context.Context, transactionID string) (Transaction, error) {
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Transaction{}, errors.WithErrorMessage(err, errors.NotKnown, "beginning get transactions db transactoion")
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
			return returning, errors.WithErrorMessage(err, errors.NotFound, fmt.Sprintf("transaction %q not found", transactionID))
		}
		return returning, errors.WithErrorMessage(err, errors.NotKnown, "scanning transaction result row")
	}

	returning.Entries, err = getEntriesByTransactionID(ctx, tx, transactionID)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return returning, errors.Wrap(err, "rolling back on fetching transaction entries error")
		}
		return returning, errors.Wrap(err, fmt.Sprintf("getting entries by transaction %q", transactionID))
	}

	if err := validate.Validate(returning); err != nil {
		return returning, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transaction fetched from database")
	}

	if err := tx.Commit(ctx); err != nil {
		return returning, errors.WithErrorMessage(err, errors.NotKnown, "committing get transaction query")
	}

	return returning, nil
}

// ListTransactions get's transactions paginated by cursor and limit
func (r *Repository) ListTransactions(ctx context.Context, transactionID string, createdAt time.Time, limit uint64) ([]Transaction, error) {
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotKnown, "beginning get transactions db transactoion")
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
		return nil, errors.WithErrorMessage(err, errors.NotKnown, "fetching transactions from database")
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
			return nil, errors.WithErrorMessage(err, errors.NotKnown, "scanning transaction result row")
		}
		returning = append(returning, transaction)
	}

	for i, transaction := range returning {
		entries, err := getEntriesByTransactionID(ctx, tx, transaction.ID)
		if err != nil {
			return nil, errors.Wrap(err, "getting entries by transaction id")
		}
		returning[i].Entries = append(returning[i].Entries, entries...)
	}

	if err := validate.Validate(returning); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transactions fetched from database")
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotKnown, "committing get transactions transaction")
	}

	return returning, nil
}

// CreateTransaction creates a transaction and associated entries in a single transaction
func (r *Repository) CreateTransaction(ctx context.Context, create Transaction) (Transaction, error) {
	// get a db-transaction
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Transaction{}, errors.WithErrorMessage(err, errors.NotKnown, "beginning create transaction db transaction")
	}
	returning := Transaction{}
	if err := tx.QueryRow(ctx,
		`INSERT INTO transaction (id, description, created_at) VALUES ($1, $2, $3)
		RETURNING id, description, created_at
		;`, create.ID, create.Description, create.CreatedAt).Scan(
		&returning.ID, &returning.Description, &returning.CreatedAt,
	); err != nil {
		return returning, errors.WithErrorMessage(err, errors.NotKnown, "scanning transaction returned from insert")
	}

	// insert the entries
	for _, entry := range create.Entries {
		returningEntry := Entry{}
		if err := tx.QueryRow(ctx,
			`INSERT INTO entry(id, transaction_id, debit_account, credit_account, amount_value, amount_currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, transaction_id, debit_account, credit_account, amount_value, amount_currency, created_at;`,
			entry.ID, returning.ID, entry.DebitAccount, entry.CreditAccount, entry.Amount.Value, entry.Amount.Currency, entry.CreatedAt).Scan(
			&returningEntry.ID,
			&returningEntry.TransactionID,
			&returningEntry.DebitAccount,
			&returningEntry.CreditAccount,
			&returningEntry.Amount.Value,
			&returningEntry.Amount.Currency,
			&returningEntry.CreatedAt,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.ForeignKeyViolation {
					return returning, errors.WithErrorMessage(err, errors.NoRelationshipFound, "inserting entry with foreign key constraint on transaction id")
				}
			}
			return returning, errors.WithErrorMessage(err, errors.NotKnown, "scanning entry result row")
		}
		returning.Entries = append(returning.Entries, returningEntry)
	}

	if err := validate.Validate(returning); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return returning, errors.WithErrorMessage(err, errors.NotKnown, "rolling back create transaction on validating transaction returned from create")
		}
		return returning, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transaction returned from creation")
	}

	if err := tx.Commit(ctx); err != nil {
		return returning, errors.WithErrorMessage(err, errors.NotKnown, "committing create transaction transaction")
	}
	return returning, nil
}
