package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	event "github.com/hampgoodwin/GoLuca/internal/event"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

func (s *Service) GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "service.GetTransaction", trace.WithAttributes(
		attribute.String("transaction_id", transactionID),
	))
	defer span.End()

	repoTransaction, err := s.repository.GetTransaction(ctx, transactionID)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("getting transaction %q: %w", transactionID, err)
	}

	transaction := transformer.NewTransactionFromRepoTransaction(repoTransaction)
	if err := validate.Validate(transaction); err != nil {
		return transaction, errors.Join(fmt.Errorf("validating transfer from repository transfer: %w", err), ierrors.ErrNotValidInternalData)
	}
	return transaction, nil
}

func (s *Service) ListTransactions(ctx context.Context, cursor string, limit uint64) ([]transaction.Transaction, string, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "service.ListTransaction", trace.WithAttributes(
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
		attribute.String("cursor.id", id),
		attribute.String("cursor.created_at", createdAt.String()),
	)

	repoTransactions, err := s.repository.ListTransactions(ctx, id, createdAt, limit)
	if err != nil {
		return nil, "", fmt.Errorf("fetching transactions from database with cursor %q", cursor, err)
	}

	transactions := []transaction.Transaction{}
	for _, repoTransaction := range repoTransactions {
		transactions = append(transactions, transformer.NewTransactionFromRepoTransaction(repoTransaction))
	}
	if err := validate.Validate(transactions); err != nil {
		return transactions, "", errors.Join(fmt.Errorf("validating transfers from repository transfers: %w", err), ierrors.ErrNotValidInternalData)
	}

	nextCursor := ""
	if len(transactions) == int(limit) {
		var err error
		lastTransaction := transactions[len(transactions)-1]
		nextCursor, err = pagination.Cursor{
			ID:         lastTransaction.ID,
			Time:       lastTransaction.CreatedAt,
			Parameters: nil, // once I add query paramters/filters, include this
		}.EncodeCursor()
		if err != nil {
			return nil, "", errors.Join(fmt.Errorf("encoding cursor for next cursor: %w", err), ierrors.ErrNotValidInternalData)
		}
		transactions = transactions[:len(transactions)-1]
	}

	return transactions, nextCursor, nil
}

func (s *Service) CreateTransaction(ctx context.Context, create transaction.Transaction) (transaction.Transaction, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.service.CreateTransaction", trace.WithAttributes(
		attribute.String("description", create.Description),
		attribute.Int64("entries", int64(len(create.Entries))),
	))
	defer span.End()

	uuidv7, err := uuid.NewV7()
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("creating uuid7: %w", err)
	}
	create.ID = uuidv7.String()
	create.CreatedAt = time.Unix(uuidv7.Time().UnixTime())
	for i := 0; i < len(create.Entries); i++ {
		_, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.service.CreateTransaction.entries")
		defer span.End()

		entryUUIDV7, err := uuid.NewV7()
		if err != nil {
			return transaction.Transaction{}, fmt.Errorf("creating uuid7 for entry: %w", err)
		}
		create.Entries[i].ID = entryUUIDV7.String()
		create.Entries[i].TransactionID = create.ID
		create.Entries[i].CreatedAt = time.Unix(uuidv7.Time().UnixTime())

		span.SetAttributes(
			attribute.String("id", create.Entries[i].ID),
			attribute.String("transaction_id", create.Entries[i].TransactionID),
			attribute.String("description", create.Entries[i].Description),
			attribute.String("debit_account", create.Entries[i].DebitAccount),
			attribute.String("credit_account", create.Entries[i].CreditAccount),
			attribute.String("amount", fmt.Sprintf("%d_%s", create.Entries[i].Amount.Value, create.Entries[i].Amount.Currency)),
		)
	}
	if err := validate.Validate(create); err != nil {
		return transaction.Transaction{}, errors.Join(fmt.Errorf("validating transaction before persisting to database: %w", err), ierrors.ErrNotValidRequestData)
	}

	repoTransaction := transformer.NewRepoTransactionFromTransaction(create)

	created, err := s.repository.CreateTransaction(ctx, repoTransaction)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("storing transaction: %w", err)
	}

	transaction := transformer.NewTransactionFromRepoTransaction(created)
	if err := validate.Validate(transaction); err != nil {
		return transaction, errors.Join(fmt.Errorf("validating transfer from repository transfer: %w", err), ierrors.ErrNotValidInternalData)
	}

	protoCreated := transformer.NewProtoTransactionFromTransaction(transaction)
	data, err := proto.Marshal(protoCreated)
	if err != nil {
		return transaction, fmt.Errorf("proto encoding created transaction: %w", err)
	}
	if err := s.publisher.Publish(event.SubjectTransactionCreated, data); err != nil {
		log.Printf("publishing %q: %v", event.SubjectTransactionCreated, err)
	}

	return transaction, nil
}
