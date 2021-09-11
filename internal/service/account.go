package service

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/data"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/account"
)

func GetAccount(ctx context.Context, acctID int64) (*account.Account, error) {
	account, err := data.GetAccount(ctx, acctID)
	if err != nil {
		return nil, errors.Wrap(err, "getting account from database")
	}
	return account, nil
}

func GetAccounts(ctx context.Context, cursor int64, limit int64) ([]account.Account, error) {
	accounts, err := data.GetAccounts(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting accounts from database")
	}
	return accounts, nil
}

func CreateAccount(ctx context.Context, account *account.Account) (*account.Account, error) {
	created, err := data.CreateAccount(ctx, account)
	if err != nil {
		return nil, errors.Wrap(err, "persisting account to database")
	}
	return created, nil
}
