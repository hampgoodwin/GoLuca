package connect

import (
	"context"
	"fmt"
	"net/http"

	accountv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1"
	accountv1connect "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1/accountv1connect"
	"github.com/hampgoodwin/GoLuca/internal/account/service"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transformer"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func Register(
	m *http.ServeMux,
	h *Handler,
) {
	// TODO: replace all method types with expectec connect types
	path, handler := accountv1connect.NewAccountServiceHandler(h)
	m.Handle(path, handler)
}

func (h *Handler) GetAccount(
	ctx context.Context,
	req *connect.Request[accountv1.GetAccountRequest],
) (*connect.Response[accountv1.GetAccountResponse], error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.GetAccount", trace.WithAttributes(
		attribute.String("account_id", req.Msg.GetAccountId()),
	))
	defer span.End()

	serviceAccount, err := h.service.GetAccount(ctx, req.Msg.AccountId)
	if err != nil {
		return nil, h.respondError(ctx, err)
	}

	account := NewProtoAccountFromAccount(serviceAccount)

	res := connect.NewResponse(&accountv1.GetAccountResponse{Account: account})

	return res, nil
}

func (h *Handler) ListAccounts(
	ctx context.Context,
	req *connect.Request[accountv1.ListAccountsRequest],
) (*connect.Response[accountv1.ListAccountsResponse], error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.ListAccounts", trace.WithAttributes(
		attribute.Int64("page_size", int64(req.Msg.GetPageSize())),
		attribute.String("page_token", req.Msg.GetPageToken()),
	))
	defer span.End()

	limit, cursor := req.Msg.PageSize, req.Msg.PageToken
	if limit == 0 {
		limit = 10
	}

	accounts, nextCursor, err := h.service.ListAccounts(ctx, cursor, limit)
	if err != nil {
		return nil, h.respondError(ctx, err)
	}

	listAccountsResponse := &accountv1.ListAccountsResponse{
		NextPageToken: nextCursor,
	}
	for _, account := range accounts {
		listAccountsResponse.Accounts = append(listAccountsResponse.Accounts, NewProtoAccountFromAccount(account))
	}

	res := connect.NewResponse(listAccountsResponse)

	return res, nil
}

func (h *Handler) CreateAccount(
	ctx context.Context,
	create *connect.Request[accountv1.CreateAccountRequest],
) (*connect.Response[accountv1.CreateAccountResponse], error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.CreateAccount", trace.WithAttributes(
		attribute.String("parent_id", create.Msg.GetParentId()),
		attribute.String("name", create.Msg.GetName()),
		attribute.String("type", create.Msg.GetType().String()),
		attribute.String("basis", create.Msg.GetBasis().String()),
	))
	defer span.End()

	serviceAccount := transformer.NewAccountFromProtoCreateAccount(create.Msg)

	serviceAccount, err := h.service.CreateAccount(ctx, serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	account := NewProtoAccountFromAccount(serviceAccount)

	res := connect.NewResponse(&accountv1.CreateAccountResponse{Account: account})

	return res, nil
}
