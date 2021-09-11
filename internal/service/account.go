package service

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/account"
)

func (s *Service) GetAccount(ctx context.Context, acctID int64) (*account.Account, error) {
	account, err := s.repository.GetAccount(ctx, acctID)
	if err != nil {
		return nil, errors.Wrap(err, "getting account from database")
	}
	return account, nil
}

func (s *Service) GetAccounts(ctx context.Context, cursor int64, limit int64) ([]account.Account, error) {
	accounts, err := s.repository.GetAccounts(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting accounts from database")
	}
	return accounts, nil
}

func (s *Service) CreateAccount(ctx context.Context, account *account.Account) (*account.Account, error) {
	created, err := s.repository.CreateAccount(ctx, account)
	if err != nil {
		return nil, errors.Wrap(err, "persisting account to database")
	}
	return created, nil
}
