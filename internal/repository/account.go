package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
	"github.com/jackc/pgx/v4"
)

// GetAccount gets an account from the database
func (r *Repository) GetAccount(ctx context.Context, accountID string) (Account, error) {
	acct := Account{}
	if err := r.database.QueryRow(ctx,
		`SELECT id, parent_id, name, type, basis, created_at
		FROM account
		WHERE id=$1
		;`,
		accountID).Scan(
		&acct.ID, &acct.ParentID, &acct.Name, &acct.Type, &acct.Basis, &acct.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return acct, errors.WithErrorMessage(err, errors.NotFound, fmt.Sprintf("account %q not found", accountID))
		}
		return acct, errors.WithErrorMessage(err, errors.NotKnown, "scanning account result row")
	}
	if err := validate.Validate(acct); err != nil {
		return acct, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating account fetched from database")
	}
	return acct, nil
}

// GetAccounts get accounts paginated based on a cursor and limit
func (r *Repository) GetAccounts(ctx context.Context, accountID string, createdAt time.Time, limit uint64) ([]Account, error) {
	query := `SELECT id, parent_id, name, type, basis, created_at
		FROM account
		WHERE 1=1`
	var params []interface{}
	if accountID != "" && !createdAt.IsZero() {
		params = append(params, accountID)
		query += fmt.Sprintf(" AND account.id <= $%d", len(params))
		params = append(params, createdAt)
		query += fmt.Sprintf(" AND account.created_at <= $%d", len(params))
	}
	query += " ORDER BY created_at DESC"
	if limit != 0 {
		params = append(params, limit)
		query += fmt.Sprintf(" LIMIT $%d", len(params))
	}
	query += ";"
	rows, err := r.database.Query(ctx, query, params...)
	if err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotKnown, "fetching accounts from data store")
	}
	defer rows.Close()
	accounts := []Account{}
	for rows.Next() {
		acct := Account{}
		if err := rows.Scan(&acct.ID, &acct.ParentID, &acct.Name, &acct.Type, &acct.Basis, &acct.CreatedAt); err != nil {
			return nil, errors.WithErrorMessage(err, errors.NotKnown, "scanning account result row")
		}
		accounts = append(accounts, acct)
	}
	if err := validate.Validate(accounts); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating accounts fetched from data store")
	}
	return accounts, nil
}

// CreateAccount creates an account record in the database and returns the created record
func (r *Repository) CreateAccount(ctx context.Context, create Account) (Account, error) {
	// get a db-transaction
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Account{}, errors.WithErrorMessage(err, errors.NotKnown, "beginning create account db transaction")
	}

	returning := Account{}
	if err := tx.QueryRow(ctx, `
		INSERT INTO account(id, parent_id, name, type, basis, created_at)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, parent_id, name, type, basis, created_at
		;`,
		create.ID, create.ParentID, create.Name, create.Type, create.Basis, create.CreatedAt).Scan(
		&returning.ID,
		&returning.ParentID,
		&returning.Name,
		&returning.Type,
		&returning.Basis,
		&returning.CreatedAt,
	); err != nil {
		return returning, errors.WithErrorMessage(err, errors.NotKnown, "scanning account returned from insert")
	}
	if err := validate.Validate(returning); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return returning, errors.WithErrorMessage(err, errors.NotValidInternalData, "rolling back transaction on failed validating account returned from insert")
		}
		return returning, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating account returned from insert")
	}
	if err := tx.Commit(ctx); err != nil {
		return returning, errors.WithErrorMessage(err, errors.NotKnown, "committing account insert transaction")
	}
	return returning, nil
}
