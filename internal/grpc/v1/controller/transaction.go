package controller

import (
	"context"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
)

func (c *Controller) GetTransaction(ctx context.Context, req *servicev1.GetTransactionRequest) (*servicev1.GetTransactionResponse, error) {
	if err := validate.Validate(req); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating request")
	}

	serviceTransaction, err := c.service.GetTransaction(ctx, req.TransactionId)
	if err != nil {
		return nil, errors.WithMessage(err, "getting account")
	}

	transaction := transformer.NewProtoTransactionFromTransaction(serviceTransaction)
	if err := validate.Validate(transaction); err != nil {
		return nil, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating transaction from service transaction")
	}

	return &servicev1.GetTransactionResponse{Transaction: transaction}, nil
}

func (c *Controller) ListTransactions(ctx context.Context, req *servicev1.ListTransactionsRequest) (*servicev1.ListTransactionsResponse, error) {
	limit, cursor := req.PageSize, req.PageToken
	if limit == 0 {
		limit = 10
	}
	if err := validate.Var(cursor, "omitempty,base64"); err != nil {
		return nil, c.respondError(errors.WithErrorMessage(err, errors.NotValidRequest, "invalid cursor or token"))
	}

	transactions, nextCursor, err := c.service.ListTransactions(ctx, cursor, limit)
	if err != nil {
		return nil, c.respondError(err)
	}

	listTransactionsResponse := &servicev1.ListTransactionsResponse{
		NextPageToken: nextCursor,
	}
	for _, transaction := range transactions {
		listTransactionsResponse.Transactions = append(listTransactionsResponse.Transactions, transformer.NewProtoTransactionFromTransaction(transaction))
	}
	if err := validate.Validate(listTransactionsResponse); err != nil {
		return nil, c.respondError(errors.WithErrorMessage(err, errors.NotValidInternalData, "validating list transactions response from transactions"))
	}

	return listTransactionsResponse, nil
}

func (c *Controller) CreateTransaction(ctx context.Context, create *servicev1.CreateTransactionRequest) (*servicev1.CreateTransactionResponse, error) {
	if err := validate.Validate(create); err != nil {
		return nil, c.respondError(errors.WithErrorMessage(err, errors.NotValidRequestData, "validating create transaction request"))
	}

	serviceTransaction := transformer.NewTransactionFromProtoCreateTransaction(create)

	createdTransaction, err := c.service.CreateTransaction(ctx, serviceTransaction)
	if err != nil {
		return nil, errors.WithMessage(err, "creating account")
	}

	transaction := transformer.NewProtoTransactionFromTransaction(createdTransaction)
	if err := validate.Validate(transaction); err != nil {
		return nil, c.respondError(errors.WithErrorMessage(err, errors.NotValidInternalData, "validating transaction from created transaction"))
	}
	return &servicev1.CreateTransactionResponse{Transaction: transaction}, nil
}
