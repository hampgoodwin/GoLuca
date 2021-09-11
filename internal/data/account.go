package data

import (
	"context"
	"database/sql"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	"github.com/jackc/pgx/v4"
)

// GetAccount gets an account from the database
func GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	account := &account.Account{}
	if err := DBPool.QueryRow(ctx, `SELECT id, parent_id, name, type, basis
FROM account WHERE id=$1;`, id).Scan(
		&account.ID, &account.ParentID, &account.Name, &account.Type, &account.Basis,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WrapFlag(err, "scanning account result row", errors.NotFound)
		}
		return nil, errors.Wrap(err, "scanning account result row")
	}
	if err := Validate(account); err != nil {
		return nil, errors.WrapFlag(err, "validating account retrieved from datastore", errors.NotValidInternalData)
	}
	return account, nil
}

// GetAccounts get accounts paginated based on a cursor and limit
func GetAccounts(ctx context.Context, cursor int64, limit int64) ([]account.Account, error) {
	rows, err := DBPool.Query(ctx, `SELECT id, parent_id, name, type, basis
FROM account
WHERE account.id > $1
LIMIT $2
;`, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying database for accounts")
	}
	defer rows.Close()
	accounts := []account.Account{}
	for rows.Next() {
		account := account.Account{}
		if err := rows.Scan(&account.ID, &account.ParentID, &account.Name, &account.Type, &account.Basis); err != nil {
			return nil, errors.Wrap(err, "scanning row from accounts query result set")
		}
		accounts = append(accounts, account)
	}
	if err := Validate(accounts); err != nil {
		return nil, errors.WrapFlag(err, "validating accounts retrieved from datastore", errors.NotValidInternalData)
	}
	return accounts, nil
}

// CreateAccount creates an account record in the database and returns the created record
func CreateAccount(ctx context.Context, acc *account.Account) (*account.Account, error) {
	// get a db-transaction
	tx, err := DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "starting db-transaction for creating transaction")
	}

	account := &account.Account{}
	if err := tx.QueryRow(ctx, `INSERT INTO account(parent_id, name, type, basis)
	VALUES($1, $2, $3, $4)
	RETURNING id, parent_id, name, type, basis;`,
		acc.ParentID, acc.Name, acc.Type, acc.Basis).Scan(
		&account.ID,
		&account.ParentID,
		&account.Name,
		&account.Type,
		&account.Basis,
	); err != nil {
		return nil, errors.Wrap(err, "scanning account creation query returning result set")
	}
	if err := Validate(account); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, errors.WrapFlag(err, "rolling back account creation transaction on invalid return data", errors.NotValidInternalData)
		}
		return nil, errors.WrapFlag(err, "validating account created in datastore", errors.NotValidInternalData)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "committing account creation")
	}
	return account, nil
}
