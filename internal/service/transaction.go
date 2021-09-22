package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (s *Service) GetTransactions(ctx context.Context, cursor, limit string) ([]transaction.Transaction, *string, error) {
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return nil, nil, errors.FlagWrap(
			err, errors.NotValidRequest,
			fmt.Sprintf("failed parsing provided limit query parameter %q", limit),
			"parsing limit query param")
	}
	limitInt++ // we always want one more than the size of the page, the extra at the end of the resultset serves as starting record for the next page
	var id string
	var createdAt time.Time
	if cursor != "" {
		id, createdAt, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, errors.FlagWrap(err, errors.NotValidRequest, err.Error(), "decoding base64 cursor")
		}
	}
	transactions, err := s.repository.GetTransactions(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting transactions from database")
	}

	encodedCursor := ""
	if len(transactions) == int(limitInt) {
		encodedCursor = pagination.EncodeCursor(transactions[len(transactions)-1].CreatedAt, transactions[len(transactions)-1].ID)
		transactions = transactions[:len(transactions)-1]
	}

	return transactions, &encodedCursor, nil
}

func (s *Service) GetTransaction(ctx context.Context, transactionID string) (*transaction.Transaction, error) {
	transaction, err := s.repository.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting transaction from database")
	}
	return transaction, nil
}

func (s *Service) GetTransactionEntries(ctx context.Context, transactionID string) ([]transaction.Entry, error) {
	entries, err := s.repository.GetEntriesByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting entries by transaction")
	}
	return entries, nil
}

func (s *Service) CreateTransactionAndEntries(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error) {
	transaction.ID = uuid.New().String()
	transaction.CreatedAt = time.Now()
	for i := 0; i < len(transaction.Entries); i++ {
		transaction.Entries[i].ID = uuid.New().String()
		transaction.Entries[i].TransactionID = transaction.ID
		transaction.Entries[i].CreatedAt = transaction.CreatedAt
	}

	if err := validate.Validate(transaction); err != nil {
		return nil, errors.FlagWrap(
			err, errors.NotValidRequestData,
			"provided request body with transaction failed validation",
			"validating transaction before persisting to database")
	}

	transaction, err := s.repository.CreateTransactionAndEntries(ctx, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "storing transaction")
	}
	return transaction, nil
}
