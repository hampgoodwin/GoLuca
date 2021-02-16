package data

import (
	"context"

	"github.com/abelgoodwin1988/GoLuca/pkg/account"
	"github.com/pkg/errors"
)

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
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to prepare account insert statement")
	}
	account := &account.Account{}
	txInsertAccountStmt.QueryRowContext(ctx, acc.ParentID, acc.Name, acc.Type, acc.Basis).Scan(
		&account.ID,
		&account.ParentID,
		&account.Name,
		&account.Type,
		&account.Basis,
	)
	if err := Validate(account); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return account, nil
}
