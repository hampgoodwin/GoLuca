package service

import (
	"context"
	"fmt"
	"time"

	event "github.com/hampgoodwin/GoLuca/internal/event"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
	"github.com/hampgoodwin/errors"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Service) GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "service.GetTransaction", trace.WithAttributes(
		attribute.String("transaction_id", transactionID),
	))
	defer span.End()

	repoTransaction, err := s.repository.GetTransaction(ctx, transactionID)
	if err != nil {
		return transaction.Transaction{}, errors.Wrap(err, fmt.Sprintf("getting transaction %q", transactionID))
	}

	transaction := transformer.NewTransactionFromRepoTransaction(repoTransaction)
	if err := validate.Validate(transaction); err != nil {
		return transaction, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transfer from repository transfer")
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
			return nil, "", errors.WithErrorMessage(err, errors.NotValidRequest, "parsing cursor object")
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
		return nil, "", errors.Wrap(err, fmt.Sprintf("fetching transactions from database with cursor %q", cursor))
	}

	transactions := []transaction.Transaction{}
	for _, repoTransaction := range repoTransactions {
		transactions = append(transactions, transformer.NewTransactionFromRepoTransaction(repoTransaction))
	}
	if err := validate.Validate(transactions); err != nil {
		return transactions, "", errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transfers from repository transfers")
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
			return nil, "", errors.WithErrorMessage(err, errors.NotValidInternalData, "encoding cursor for next cursor")
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

	create.ID = ksuid.New().String()
	create.CreatedAt = time.Now()
	for i := 0; i < len(create.Entries); i++ {
		create.Entries[i].ID = ksuid.New().String()
		create.Entries[i].TransactionID = create.ID
		create.Entries[i].CreatedAt = create.CreatedAt
		_, span := otel.Tracer(meta.ServiceName).Start(ctx, "intenal.service.CreateTransaction.entries")
		defer span.End()

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
		return transaction.Transaction{}, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating transaction before persisting to database")
	}

	repoTransaction := transformer.NewRepoTransactionFromTransaction(create)

	created, err := s.repository.CreateTransaction(ctx, repoTransaction)
	if err != nil {
		return transaction.Transaction{}, errors.Wrap(err, "storing transaction")
	}

	transaction := transformer.NewTransactionFromRepoTransaction(created)
	if err := validate.Validate(transaction); err != nil {
		return transaction, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transfer from repository transfer")
	}

	protoCreated := transformer.NewProtoTransactionFromTransaction(transaction)
	data, err := proto.Marshal(protoCreated)
	if err != nil {
		s.log.Error("proto encoding created transaction", zap.Any("proto_transaction", protoCreated))
		return transaction, errors.Wrap(err, "proto encoding created transaction")
	}
	if err := s.publisher.Publish(event.SubjectTransactionCreated, data); err != nil {
		// should we return/fail on a failed production of a msg..?
		s.log.Error("publishing transaction created message", zap.Error(err))
	}

	return transaction, nil
}
