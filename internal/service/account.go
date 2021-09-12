package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/account"
)

func (s *Service) GetAccount(ctx context.Context, accountID string) (*account.Account, error) {
	account, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "getting account from database")
	}
	return account, nil
}

func (s *Service) GetAccounts(ctx context.Context, cursor string, limit uint64) ([]account.Account, error) {
	accounts, err := s.repository.GetAccounts(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting accounts from database")
	}
	return accounts, nil
}

func (s *Service) CreateAccount(ctx context.Context, account *account.Account) (*account.Account, error) {
	account.ID = uuid.New().String()
	account.CreatedAt = time.Now().UTC()
	if err := validate.Validate(account); err != nil {
		return nil, errors.WrapFlag(err, "validating account before persisting to database", errors.NotValidRequestData)
	}

	created, err := s.repository.CreateAccount(ctx, account)
	if err != nil {
		return nil, errors.Wrap(err, "persisting account to database")
	}
	return created, nil
}
