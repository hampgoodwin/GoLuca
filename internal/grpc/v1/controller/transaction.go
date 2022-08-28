package controller

import (
	"context"
	"fmt"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/errors"
)

func (c *Controller) GetTransaction(ctx context.Context, req *servicev1.GetTransactionRequest) (*servicev1.GetTransactionResponse, error) {
	c.log.Info(fmt.Sprintf("getting transaction %q", req.TransactionId)) // TODO factor out in a gRPC router logger for incoming/outgoing requests.

	serviceTransaction, err := c.service.GetTransaction(ctx, req.TransactionId)
	if err != nil {
		return nil, errors.WithMessage(err, "getting account")
	}
	transaction := transformer.NewPBTransactionFromTransaction(serviceTransaction)
	return &servicev1.GetTransactionResponse{Transaction: transaction}, nil
}

func (c *Controller) ListTransactions(ctx context.Context, req *servicev1.ListTransactionsRequest) (*servicev1.ListTransactionsResponse, error) {
	return &servicev1.ListTransactionsResponse{}, nil
}

func (c *Controller) CreateTransaction(ctx context.Context, create *servicev1.CreateTransactionRequest) (*servicev1.CreateTransactionResponse, error) {

	return &servicev1.CreateTransactionResponse{}, nil
}
