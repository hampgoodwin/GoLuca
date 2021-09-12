package repository

import (
	"context"
	"database/sql"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
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
		if errors.Is(err, sql.ErrNoRows) {
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
func (r *Repository) GetTransactions(ctx context.Context, cursor string, limit uint64) ([]transaction.Transaction, error) {
	rows, err := r.Database.Query(ctx,
		`SELECT id, description
		FROM transaction
		WHERE transaction.id > $1
		LIMIT $2
		;`, cursor, limit)
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
			`INSERT INTO entry(id, transaction_id, account_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)
			RETURNING id, transaction_id, account_id, amount, created_at;`,
			entry.ID, returningTransaction.ID, entry.AccountID, entry.Amount, entry.CreatedAt).Scan(
			&returningEntry.ID,
			&returningEntry.TransactionID,
			&returningEntry.AccountID,
			&returningEntry.Amount,
			&returningEntry.CreatedAt,
		); err != nil {
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
