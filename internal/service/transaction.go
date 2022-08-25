package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/pagination"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/hampgoodwin/errors"
)

func (s *Service) GetTransactions(ctx context.Context, cursor, limit string) ([]transaction.Transaction, *string, error) {
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
	transactions, err := s.repository.GetTransactions(ctx, id, createdAt, limitInt)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("fetching transactions from database with cursor %q", cursor))
	}

	encodedCursor := ""
	if len(transactions) == int(limitInt) {
		encodedCursor = pagination.EncodeCursor(transactions[len(transactions)-1].CreatedAt, transactions[len(transactions)-1].ID)
		transactions = transactions[:len(transactions)-1]
	}

	return transactions, &encodedCursor, nil
}

func (s *Service) GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	txn, err := s.repository.GetTransaction(ctx, transactionID)
	if err != nil {
		return transaction.Transaction{}, errors.Wrap(err, fmt.Sprintf("getting transaction %q", transactionID))
	}
	return txn, nil
}

func (s *Service) CreateTransactionAndEntries(ctx context.Context, txn transaction.Transaction) (transaction.Transaction, error) {
	txn.ID = uuid.New().String()
	txn.CreatedAt = time.Now()
	for i := 0; i < len(txn.Entries); i++ {
		txn.Entries[i].ID = uuid.New().String()
		txn.Entries[i].TransactionID = txn.ID
		txn.Entries[i].CreatedAt = txn.CreatedAt
	}

	if err := validate.Validate(txn); err != nil {
		return transaction.Transaction{}, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating transaction before persisting to database")
	}

	txn, err := s.repository.CreateTransaction(ctx, txn)
	if err != nil {
		return transaction.Transaction{}, errors.Wrap(err, "storing transaction")
	}
	return txn, nil
}
