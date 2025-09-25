package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	transactionv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/transaction/v1"
	"github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/transaction/v1/transactionv1connect"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transaction/service"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Register(
	m *http.ServeMux,
	h *Handler,
) {
	// TODO: replace all method types with expectec connect types
	path, handler := transactionv1connect.NewTransactionServiceHandler(h)
	m.Handle(path, handler)
}

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (c *Handler) GetTransaction(
	ctx context.Context,
	req *connect.Request[transactionv1.GetTransactionRequest],
) (*connect.Response[transactionv1.GetTransactionResponse], error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.GetTransaction", trace.WithAttributes(
		attribute.String("transaction_id", req.Msg.GetTransactionId()),
	))
	defer span.End()

	if err := validate.Validate(req); err != nil {
		return nil, errors.Join(fmt.Errorf("validating request: %w", err), ierrors.ErrNotValidRequestData)
	}

	serviceTransaction, err := c.service.GetTransaction(ctx, req.Msg.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	transaction := NewProtoTransactionFromTransaction(serviceTransaction)
	if err := validate.Validate(transaction); err != nil {
		return nil, errors.Join(fmt.Errorf("validating transaction from service transaction: %w", err), ierrors.ErrNotValidRequestData)
	}

	res := connect.NewResponse(&transactionv1.GetTransactionResponse{Transaction: transaction})
	return res, nil
}

func (c *Handler) ListTransactions(
	ctx context.Context,
	req *connect.Request[transactionv1.ListTransactionsRequest],
) (*connect.Response[transactionv1.ListTransactionsResponse], error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.ListTransaction", trace.WithAttributes(
		attribute.Int64("page_size", int64(req.Msg.GetPageSize())),
		attribute.String("page_token", req.Msg.GetPageToken()),
	))
	defer span.End()

	limit, cursor := req.Msg.PageSize, req.Msg.PageToken
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

	listTransactionsResponse := &transactionv1.ListTransactionsResponse{
		NextPageToken: nextCursor,
	}
	for _, transaction := range transactions {
		listTransactionsResponse.Transactions = append(listTransactionsResponse.Transactions, NewProtoTransactionFromTransaction(transaction))
	}
	if err := validate.Validate(listTransactionsResponse); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating list transactions response from transactions: %w", err), ierrors.ErrNotValidInternalData))
	}

	res := connect.NewResponse(listTransactionsResponse)
	return res, nil
}

func (c *Handler) CreateTransaction(
	ctx context.Context,
	create *connect.Request[transactionv1.CreateTransactionRequest],
) (*connect.Response[transactionv1.CreateTransactionResponse], error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.CreateTransaction", trace.WithAttributes(
		attribute.String("parent_id", create.Msg.GetDescription()),
		attribute.Int64("count_entries", int64(len(create.Msg.GetEntries()))),
	))
	defer span.End()

	if err := validate.Validate(create); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating create transaction request: %w", err), ierrors.ErrNotValidRequestData))
	}

	serviceTransaction := NewTransactionFromProtoCreateTransaction(create.Msg)

	createdTransaction, err := c.service.CreateTransaction(ctx, serviceTransaction)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	transaction := NewProtoTransactionFromTransaction(createdTransaction)
	if err := validate.Validate(transaction); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating transaction from created transaction: %w", err), ierrors.ErrNotValidInternalData))
	}

	res := connect.NewResponse(&transactionv1.CreateTransactionResponse{Transaction: transaction})
	return res, nil
}
