package connect

import (
	"context"
	"fmt"
	"net/http"

	transactionv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/transaction/v1"
	"github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/transaction/v1/transactionv1connect"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transaction"
	"github.com/hampgoodwin/GoLuca/internal/transaction/service"

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

	serviceTransaction, err := c.service.GetTransaction(ctx, req.Msg.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	txn := transaction.NewProtoTransactionFromTransaction(serviceTransaction)

	res := connect.NewResponse(&transactionv1.GetTransactionResponse{Transaction: txn})
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

	txn, nextCursor, err := c.service.ListTransactions(ctx, cursor, limit)
	if err != nil {
		return nil, c.respondError(ctx, err)
	}

	listTransactionsResponse := &transactionv1.ListTransactionsResponse{
		NextPageToken: nextCursor,
	}
	for _, txn := range txn {
		listTransactionsResponse.Transactions = append(listTransactionsResponse.Transactions, transaction.NewProtoTransactionFromTransaction(txn))
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

	serviceTransaction := transaction.NewTransactionFromProtoCreateTransaction(create.Msg)

	createdTransaction, err := c.service.CreateTransaction(ctx, serviceTransaction)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	transaction := transaction.NewProtoTransactionFromTransaction(createdTransaction)

	res := connect.NewResponse(&transactionv1.CreateTransactionResponse{Transaction: transaction})
	return res, nil
}
