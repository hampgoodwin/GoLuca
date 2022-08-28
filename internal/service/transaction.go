package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/transaction"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
	"github.com/hampgoodwin/errors"
	"github.com/segmentio/ksuid"
)

func (s *Service) GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error) {
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

func (s *Service) ListTransactions(ctx context.Context, cursor, limit string) ([]transaction.Transaction, *string, error) {
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parsing limit query parameter")
	}
	limitInt++ // we always want one more than the size of the page, the extra at the end of the resultset serves as starting record for the next page

	var id string
	var createdAt time.Time
	if cursor != "" {
		cursor, err := pagination.ParseCursor(cursor)
		if err != nil {
			return nil, nil, errors.WithErrorMessage(err, errors.NotValidRequest, "parsing cursor object")
		}
		id = cursor.ID
		createdAt = cursor.Time
	}

	repoTransactions, err := s.repository.ListTransactions(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("fetching transactions from database with cursor %q", cursor))
	}

	transactions := []transaction.Transaction{}
	for _, repoTransaction := range repoTransactions {
		transactions = append(transactions, transformer.NewTransactionFromRepoTransaction(repoTransaction))
	}
	if err := validate.Validate(transactions); err != nil {
		return transactions, nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transfers from repository transfers")
	}

	nextCursor := ""
	if len(transactions) == int(limitInt) {
		var err error
		lastTransaction := transactions[len(transactions)-1]
		nextCursor, err = pagination.Cursor{
			ID:         lastTransaction.ID,
			Time:       lastTransaction.CreatedAt,
			Parameters: nil, // once I add query paramters/filters, include this
		}.EncodeCursor()
		if err != nil {
			return nil, nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "encoding cursor for next cursor")
		}
		transactions = transactions[:len(transactions)-1]
	}

	return transactions, &nextCursor, nil
}

func (s *Service) CreateTransaction(ctx context.Context, create transaction.Transaction) (transaction.Transaction, error) {
	create.ID = ksuid.New().String()
	create.CreatedAt = time.Now()
	for i := 0; i < len(create.Entries); i++ {
		create.Entries[i].ID = ksuid.New().String()
		create.Entries[i].TransactionID = create.ID
		create.Entries[i].CreatedAt = create.CreatedAt
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

	return transaction, nil
}
