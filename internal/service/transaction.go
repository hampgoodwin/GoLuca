package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (s *Service) GetTransactions(ctx context.Context, cursor string, limit uint64) ([]transaction.Transaction, error) {
	transactions, err := s.repository.GetTransactions(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting transactions from database")
	}
	return transactions, nil
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
		return nil, errors.WrapFlag(err, "validating transaction before persisting to database", errors.NotValidRequestData)
	}

	transaction, err := s.repository.CreateTransactionAndEntries(ctx, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "storing transaction")
	}
	return transaction, nil
}
