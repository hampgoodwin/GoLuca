package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
)

func (s *Service) GetAccount(ctx context.Context, accountID string) (*account.Account, error) {
	account, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) GetAccounts(ctx context.Context, cursor, limit string) ([]account.Account, *string, error) {
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return nil, nil, err
	}
	limitInt++ // we always want one more than the size of the page, the extra at the end of the resultset serves as starting record for the next page
	var id string
	var createdAt time.Time
	if cursor != "" {
		id, createdAt, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, err
		}
	}
	accounts, err := s.repository.GetAccounts(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, err
	}

	nextCursor := ""
	if len(accounts) == int(limitInt) {
		nextCursor = pagination.EncodeCursor(accounts[len(accounts)-1].CreatedAt, accounts[len(accounts)-1].ID)
		accounts = accounts[:len(accounts)-1]
	}
	return accounts, &nextCursor, nil
}

func (s *Service) CreateAccount(ctx context.Context, account *account.Account) (*account.Account, error) {
	account.ID = uuid.New().String()
	account.CreatedAt = time.Now()

	if err := validate.Validate(account); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return created, nil
}
