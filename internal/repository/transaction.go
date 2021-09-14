package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

// GetTransaction get's a transaction record, without it's entries, by the transaction ID
func (r *Repository) GetTransaction(ctx context.Context, transactionID string) (*transaction.Transaction, error) {
	returningTransaction := &transaction.Transaction{}
	if err := r.Database.QueryRow(ctx,
		`SELECT id, description, created_at
		FROM transaction
		WHERE id=$1
		;`, transactionID).Scan(
		&returningTransaction.ID,
		&returningTransaction.Description,
		&returningTransaction.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WrapFlag(err, "scanning transaction result row", errors.NotFound)
		}
		return nil, errors.Wrap(err, "scanning row from transaction query result set")
	}

	var err error
	returningTransaction.Entries, err = r.GetEntriesByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "querying databse for entries while getting transaction")
	}

	if err := validate.Validate(returningTransaction); err != nil {
		return nil, errors.WrapFlag(err, "validating transactions retrieved from database", errors.NotValidInternalData)
	}

	return returningTransaction, nil
}

// GetTransactions get's transactions paginaged by cursor and limit
func (r *Repository) GetTransactions(ctx context.Context, transactionID string, createdAt time.Time, limit uint64) ([]transaction.Transaction, error) {
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

	rows, err := r.Database.Query(ctx, query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "querying database for transactions")
	}
	defer rows.Close()
	returningTransactions := []transaction.Transaction{}
	for rows.Next() {
		transaction := transaction.Transaction{}
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.CreatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning transactions results set")
		}
		returningTransactions = append(returningTransactions, transaction)
	}

	for i, transaction := range returningTransactions {
		entries, err := r.GetEntriesByTransactionID(ctx, transaction.ID)
		if err != nil {
			return nil, errors.Wrap(err, "querying databse for entries while getting transactions")
		}
		returningTransactions[i].Entries = append(returningTransactions[i].Entries, entries...)
	}

	if err := validate.Validate(returningTransactions); err != nil {
		return nil, errors.WrapFlag(err, "validating transactions retrieved from database", errors.NotValidInternalData)
	}

	return returningTransactions, nil
}

// CreateTransactionAndEntries creates a transaction and associated entries in a single transaction
func (r *Repository) CreateTransactionAndEntries(ctx context.Context, create *transaction.Transaction) (*transaction.Transaction, error) {
	// get a db-transaction
	tx, err := r.Database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "starting database transaction for creating transaction")
	}
	returningTransaction := &transaction.Transaction{}
	if err := tx.QueryRow(ctx,
		`INSERT INTO transaction (id, description, created_at) VALUES ($1, $2, $3)
		RETURNING id, description, created_at
		;`, create.ID, create.Description, create.CreatedAt).Scan(
		&returningTransaction.ID, &returningTransaction.Description, &returningTransaction.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "scanning transaction created return result set")
	}

	// insert the entries
	for _, entry := range create.Entries {
		returningEntry := transaction.Entry{}
		if err := tx.QueryRow(ctx,
			`INSERT INTO entry(id, transaction_id, debit_account, credit_account, amount_value, amount_currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, transaction_id, debit_account, credit_account, amount_value, amount_currency, created_at;`,
			entry.ID, returningTransaction.ID, entry.DebitAccount, entry.CreditAccount, entry.Amount.Value, entry.Amount.Currency, entry.CreatedAt).Scan(
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
					return nil, errors.WrapFlag(err, "scanning transaction entry created return result set", errors.NoRelationshipFound)
				}
			}
			return nil, errors.Wrap(err, "scanning transaction entry created return result set")
		}
		returningTransaction.Entries = append(returningTransaction.Entries, returningEntry)
	}

	if err := validate.Validate(returningTransaction); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, errors.WrapFlag(err, "rolling back transaction creation db-transaction on invalid return data", errors.NotValidInternalData)
		}
		return nil, errors.WrapFlag(err, "validating transaction created in datastore", errors.NotValidInternalData)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "committing transaction creation")
	}
	return returningTransaction, nil
}
