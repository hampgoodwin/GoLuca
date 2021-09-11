package service

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (s *Service) GetTransactions(ctx context.Context, cursor int64, limit int64) ([]transaction.Transaction, error) {
	transactions, err := s.repository.GetTransactions(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting transactions from database")
	}
	return transactions, nil
}

func (s *Service) GetTransaction(ctx context.Context, transactionID int64) (*transaction.Transaction, error) {
	transaction, err := s.repository.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting transaction from database")
	}
	return transaction, nil
}

func (s *Service) GetTransactionEntries(ctx context.Context, transactionID int64) ([]transaction.Entry, error) {
	entries, err := s.repository.GetEntriesByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting entries by transaction")
	}
	return entries, nil
}

func (s *Service) CreateTransaction(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error) {
	transaction, err := s.repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "storing transaction")
	}
	return transaction, nil
}
