package data

import (
	"context"

	"github.com/abelgoodwin1988/GoLuca/pkg/account"
	"github.com/pkg/errors"
)

// GetAccount gets an account from the database
func GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	selectAccountStmt, err := DB.PrepareContext(ctx, `SELECT id, parent_id, name, type, basis
FROM account WHERE id=$1;`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare account select statement")
	}
	account := &account.Account{}
	if err := selectAccountStmt.QueryRowContext(ctx, id).
		Scan(&account.ID, &account.ParentID, &account.Name, &account.Type, &account.Basis); err != nil {
		return nil, errors.Wrap(err, "failed to scan row from account query results set")
	}
	if err := Validate(account); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccounts get accounts paginated based on a cursor and limit
func GetAccounts(ctx context.Context, cursor int64, limit int64) ([]account.Account, error) {
	accountsSelectStmt, err := DB.PrepareContext(ctx, `
SELECT id, parent_id, name, type, basis
FROM account
WHERE account.id => $1
LIMIT $2
;`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare account select statement")
	}
	rows, err := accountsSelectStmt.QueryContext(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "error quering database for accounts")
	}
	defer rows.Close()
	accounts := []account.Account{}
	for rows.Next() {
		account := account.Account{}
		if err := rows.Scan(&account.ID, &account.ParentID, &account.Name, &account.Type, &account.Basis); err != nil {
			return nil, errors.Wrap(err, "failed to scan row from accounts query result set")
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// CreateAccount creates an account record in the database and returns the created record
func CreateAccount(ctx context.Context, acc *account.Account) (*account.Account, error) {
	// get a db-transaction
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start db-transaction for creating transaction")
	}

	txInsertAccountStmt, err := tx.PrepareContext(ctx, `INSERT INTO account(parent_id, name, type, basis)
	VALUES($1, $2, $3, $4)
	RETURNING id, parent_id, name, type, basis;`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, errors.Wrap(err, "failed to roll back account creation transaction")
		}
		return nil, errors.Wrap(err, "failed to prepare account insert statement")
	}
	account := &account.Account{}
	if err := txInsertAccountStmt.QueryRowContext(ctx, acc.ParentID, acc.Name, acc.Type, acc.Basis).Scan(
		&account.ID,
		&account.ParentID,
		&account.Name,
		&account.Type,
		&account.Basis,
	); err != nil {
		return nil, errors.Wrap(err, "failed to scan account creation query result set")
	}
	if err := Validate(account); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, errors.Wrap(err, "failed to rollback account creation transaction on failed return data validation")
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, errors.Wrap(err, "failed to rollback back on account creation failed commit")
		}
		return nil, errors.Wrap(err, "failed to commit on account creation")
	}
	return account, nil
}
