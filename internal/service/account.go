package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/account"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
	"github.com/hampgoodwin/errors"
	"github.com/segmentio/ksuid"
)

func (s *Service) GetAccount(ctx context.Context, accountID string) (account.Account, error) {
	repoAccount, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return account.Account{}, errors.Wrap(err, "fetching account from database")
	}

	account := transformer.NewAccountFromRepoAccount(repoAccount)
	if err := validate.Validate(account); err != nil {
		return account, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating account from repository account")
	}

	return account, nil
}

func (s *Service) ListAccounts(ctx context.Context, cursor, limit string) ([]account.Account, *string, error) {
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parsing limit query parameter")
	}
	limitInt++ // we always want one more than the size of the page, the extra at the end of the resultset serves as starting record for the next page
	var id string
	var createdAt time.Time
	if cursor != "" {
		id, createdAt, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, errors.Wrap(errors.WithMessage(err, err.Error()), "decoding cursor")
			// return nil, nil, errors.WrapWithErrorMessage(err, errors.NotValidRequest, err.Error(), "decoding cursor")
		}
	}
	repoAccounts, err := s.repository.ListAccounts(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("fetching accounts from database with cursor %q", cursor))
	}

	nextCursor := ""
	if len(repoAccounts) == int(limitInt) {
		nextCursor = pagination.EncodeCursor(repoAccounts[len(repoAccounts)-1].CreatedAt, repoAccounts[len(repoAccounts)-1].ID)
		repoAccounts = repoAccounts[:len(repoAccounts)-1]
	}

	accounts := []account.Account{}
	for _, repoAccount := range repoAccounts {
		accounts = append(accounts, transformer.NewAccountFromRepoAccount(repoAccount))
	}
	if err := validate.Validate(accounts); err != nil {
		return accounts, nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating accounts from repository accounts")
	}

	return accounts, &nextCursor, nil
}

func (s *Service) CreateAccount(ctx context.Context, create account.Account) (account.Account, error) {
	create.ID = ksuid.New().String()
	create.CreatedAt = time.Now()

	if err := validate.Validate(create); err != nil {
		return account.Account{}, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating account")
	}

	repoAccount := transformer.NewRepoAccountFromAccount(create)

	created, err := s.repository.CreateAccount(ctx, repoAccount)
	if err != nil {
		return account.Account{}, errors.Wrap(err, "creating account in database")
	}

	account := transformer.NewAccountFromRepoAccount(created)
	if err := validate.Validate(account); err != nil {
		return account, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating account from repository account")
	}

	return account, nil
}
