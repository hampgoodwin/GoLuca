package service

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/data"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/pkg/errors"
)

func GetTransactions(ctx context.Context, cursor int64, limit int64) ([]transaction.Transaction, error) {
	transactions, err := data.GetTransactions(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting transactions from database")
	}
	return transactions, nil
}

func GetTransaction(ctx context.Context, transactionID int64) (*transaction.Transaction, error) {
	transaction, err := data.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting transaction from database")
	}
	return transaction, nil
}

func GetTransactionEntries(ctx context.Context, transactionID int64) ([]transaction.Entry, error) {
	entries, err := data.GetEntriesByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "getting entries by transaction")
	}
	return entries, nil
}

func CreateTransaction(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error) {
	transaction, err := data.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "storing transaction")
	}
	return transaction, nil
}
