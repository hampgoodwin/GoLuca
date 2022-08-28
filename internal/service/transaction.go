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
		id, createdAt, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, errors.WrapWithErrorMessage(err, errors.NotKnown, err.Error(), "decoding cursor")
		}
	}
	repoTransactions, err := s.repository.ListTransactions(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("fetching transactions from database with cursor %q", cursor))
	}

	encodedCursor := ""
	if len(repoTransactions) == int(limitInt) {
		encodedCursor = pagination.EncodeCursor(repoTransactions[len(repoTransactions)-1].CreatedAt, repoTransactions[len(repoTransactions)-1].ID)
		repoTransactions = repoTransactions[:len(repoTransactions)-1]
	}

	transactions := []transaction.Transaction{}
	for _, repoTransaction := range repoTransactions {
		transactions = append(transactions, transformer.NewTransactionFromRepoTransaction(repoTransaction))
	}
	if err := validate.Validate(transactions); err != nil {
		return transactions, nil, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transfers from repository transfers")
	}

	return transactions, &encodedCursor, nil
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
