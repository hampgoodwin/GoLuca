package controller

import (
	"context"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
)

func (c *Controller) GetTransaction(ctx context.Context, req *servicev1.GetTransactionRequest) (*servicev1.GetTransactionResponse, error) {
	return &servicev1.GetTransactionResponse{}, nil
}

func (c *Controller) ListTransactions(ctx context.Context, req *servicev1.ListTransactionsRequest) (*servicev1.ListTransactionsResponse, error) {
	return &servicev1.ListTransactionsResponse{}, nil
}

func (c *Controller) CreateTransaction(ctx context.Context, create *servicev1.CreateTransactionRequest) (*servicev1.CreateTransactionResponse, error) {
	return &servicev1.CreateTransactionResponse{}, nil
}
