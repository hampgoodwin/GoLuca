package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/hampgoodwin/GoLuca/internal/account"
	"github.com/hampgoodwin/GoLuca/internal/account/event"
	"github.com/hampgoodwin/GoLuca/internal/account/repository"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
)

type Service struct {
	repository *repository.Repository
	publisher  *nats.Conn
}

func NewService(repository *repository.Repository, nec *nats.Conn) *Service {
	return &Service{repository, nec}
}

func (s *Service) GetAccount(ctx context.Context, accountID string) (account.Account, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "service.GetAccount", trace.WithAttributes(
		attribute.String("account_id", accountID),
	))
	defer span.End()

	repoAccount, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return account.Account{}, fmt.Errorf("fetching account from database: %w", err)
	}

	account := newAccountFromRepoAccount(repoAccount)

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
			return nil, "", errors.Join(fmt.Errorf("parsing cursor object: %w", err), ierrors.ErrNotValidRequest)
		}
		id = cursor.ID
		createdAt = cursor.Time
	}
	span.SetAttributes(
		attribute.String("cursor_id", id),
		attribute.String("cursor_created_at", createdAt.String()))

	repoAccounts, err := s.repository.ListAccounts(ctx, id, createdAt, limit)
	if err != nil {
		return nil, "", fmt.Errorf("fetching accounts from database with cursor %q: %w", cursor, err)
	}

	accounts := []account.Account{}
	for _, repoAccount := range repoAccounts {
		accounts = append(accounts, newAccountFromRepoAccount(repoAccount))
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
			return nil, "", errors.Join(fmt.Errorf("encoding cursor for next cursor: %w", err), ierrors.ErrNotValidInternalData)
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

	uuidv7, err := uuid.NewV7()
	if err != nil {
		return account.Account{}, fmt.Errorf("creating uuid7: %w", err)
	}
	create.ID = uuidv7.String()
	create.CreatedAt = time.Unix(uuidv7.Time().UnixTime())
	span.SetAttributes(attribute.String("id", create.ID))

	repoAccount := newRepoAccountFromAccount(create)

	created, err := s.repository.CreateAccount(ctx, repoAccount)
	if err != nil {
		return account.Account{}, fmt.Errorf("creating account in database: %w", err)
	}

	account := newAccountFromRepoAccount(created)

	protoCreated := event.NewProtoAccountFromAccount(account)
	data, err := proto.Marshal(protoCreated)
	if err != nil {
		return account, fmt.Errorf("proto encoding created account: %w", err)
	}
	if err := s.publisher.Publish(event.SubjectAccountCreated, data); err != nil {
		log.Printf("publishing %q: %v", event.SubjectAccountCreated, err)
	}

	return account, nil
}
