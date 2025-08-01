package service

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/hampgoodwin/errors"

	"github.com/hampgoodwin/GoLuca/internal/account"
	event "github.com/hampgoodwin/GoLuca/internal/event"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
)

func (s *Service) GetAccount(ctx context.Context, accountID string) (account.Account, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "service.GetAccount", trace.WithAttributes(
		attribute.String("account_id", accountID),
	))
	defer span.End()

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

func (s *Service) ListAccounts(ctx context.Context, cursor string, limit uint64) ([]account.Account, string, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.service.CreateAccount", trace.WithAttributes(
		attribute.String("cursor", cursor),
		attribute.Int64("limit", int64(limit)),
	))
	defer span.End()
	limit++ // we always want one more than the size of the page, the extra at the end of the resultset serves as starting record for the next page

	var id string
	var createdAt time.Time
	if cursor != "" {
		cursor, err := pagination.ParseCursor(cursor)
		if err != nil {
			if err := validate.Validate(cursor); err != nil {
				return nil, "", errors.WithErrorMessage(err, errors.NotValidRequestData, "invalid cursor/token")
			}
			return nil, "", errors.WithErrorMessage(err, errors.NotValidRequest, "parsing cursor object")
		}
		id = cursor.ID
		createdAt = cursor.Time
	}
	span.SetAttributes(
		attribute.String("cursor_id", id),
		attribute.String("cursor_created_at", createdAt.String()))

	repoAccounts, err := s.repository.ListAccounts(ctx, id, createdAt, limit)
	if err != nil {
		return nil, "", errors.Wrap(err, fmt.Sprintf("fetching accounts from database with cursor %q", cursor))
	}

	accounts := []account.Account{}
	for _, repoAccount := range repoAccounts {
		accounts = append(accounts, transformer.NewAccountFromRepoAccount(repoAccount))
	}
	if err := validate.Validate(accounts); err != nil {
		return accounts, "", errors.WithErrorMessage(err, errors.NotValidInternalData, "validating accounts from repository accounts")
	}

	nextCursor := ""
	if len(accounts) == int(limit) {
		var err error
		lastAccount := accounts[len(accounts)-1]
		nextCursor, err = pagination.Cursor{
			ID:         lastAccount.ID,
			Time:       lastAccount.CreatedAt,
			Parameters: map[string][]string{"previous_cursor": {cursor}}, // once I add query paramters/filters, include this
		}.EncodeCursor()
		if err != nil {
			return nil, "", errors.WithErrorMessage(err, errors.NotValidInternalData, "encoding cursor for next cursor")
		}
		accounts = accounts[:len(accounts)-1]
	}

	return accounts, nextCursor, nil
}

func (s *Service) CreateAccount(ctx context.Context, create account.Account) (account.Account, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.service.CreateAccount", trace.WithAttributes(
		attribute.String("parent_id", create.ParentID),
		attribute.String("name", create.Name),
		attribute.String("type", create.Type.String()),
		attribute.String("basis", create.Basis.String()),
	))
	defer span.End()

	create.ID = ksuid.New().String()
	create.CreatedAt = time.Now()
	span.SetAttributes(attribute.String("id", create.ID))

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

	protoCreated := transformer.NewProtoAccountFromAccount(account)
	data, err := proto.Marshal(protoCreated)
	if err != nil {
		s.log.Error("proto encoding created account", zap.Any("account", protoCreated))
		return account, errors.Wrap(err, "proto encoding created account")
	}
	if err := s.publisher.Publish(event.SubjectAccountCreated, data); err != nil {
		// should we return/fail on a failed production of a msg..?
		s.log.Error("publishing account created message", zap.Error(err))
	}

	return account, nil
}
