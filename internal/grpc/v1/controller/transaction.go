package controller

import (
	"context"
	"errors"
	"fmt"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *Controller) GetTransaction(ctx context.Context, req *servicev1.GetTransactionRequest) (*servicev1.GetTransactionResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.GetTransaction", trace.WithAttributes(
		attribute.String("transaction_id", req.GetTransactionId()),
	))
	defer span.End()

	if err := validate.Validate(req); err != nil {
		return nil, errors.Join(fmt.Errorf("validating request: %w", err), ierrors.ErrNotValidRequestData)
	}

	serviceTransaction, err := c.service.GetTransaction(ctx, req.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	transaction := transformer.NewProtoTransactionFromTransaction(serviceTransaction)
	if err := validate.Validate(transaction); err != nil {
		return nil, errors.Join(fmt.Errorf("validating transaction from service transaction: %w", err), ierrors.ErrNotValidRequestData)
	}

	return &servicev1.GetTransactionResponse{Transaction: transaction}, nil
}

func (c *Controller) ListTransactions(ctx context.Context, req *servicev1.ListTransactionsRequest) (*servicev1.ListTransactionsResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.ListTransaction", trace.WithAttributes(
		attribute.Int64("page_size", int64(req.GetPageSize())),
		attribute.String("page_token", req.GetPageToken()),
	))
	defer span.End()

	limit, cursor := req.PageSize, req.PageToken
	if limit == 0 {
		limit = 10
	}
	if err := validate.Var(cursor, "omitempty,base64"); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("invalid cursor or token: %w", err), ierrors.ErrNotValidRequest))
	}

	transactions, nextCursor, err := c.service.ListTransactions(ctx, cursor, limit)
	if err != nil {
		return nil, c.respondError(ctx, err)
	}

	listTransactionsResponse := &servicev1.ListTransactionsResponse{
		NextPageToken: nextCursor,
	}
	for _, transaction := range transactions {
		listTransactionsResponse.Transactions = append(listTransactionsResponse.Transactions, transformer.NewProtoTransactionFromTransaction(transaction))
	}
	if err := validate.Validate(listTransactionsResponse); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating list transactions response from transactions: %w", err), ierrors.ErrNotValidInternalData))
	}

	return listTransactionsResponse, nil
}

func (c *Controller) CreateTransaction(ctx context.Context, create *servicev1.CreateTransactionRequest) (*servicev1.CreateTransactionResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.CreateTransaction", trace.WithAttributes(
		attribute.String("parent_id", create.GetDescription()),
		attribute.Int64("count_entries", int64(len(create.GetEntries()))),
	))
	defer span.End()

	if err := validate.Validate(create); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating create transaction request: %w", err), ierrors.ErrNotValidRequestData))
	}

	serviceTransaction := transformer.NewTransactionFromProtoCreateTransaction(create)

	createdTransaction, err := c.service.CreateTransaction(ctx, serviceTransaction)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	transaction := transformer.NewProtoTransactionFromTransaction(createdTransaction)
	if err := validate.Validate(transaction); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating transaction from created transaction: %w", err), ierrors.ErrNotValidInternalData))
	}
	return &servicev1.CreateTransactionResponse{Transaction: transaction}, nil
}
