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

func (c *Controller) GetAccount(ctx context.Context, req *servicev1.GetAccountRequest) (*servicev1.GetAccountResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.GetAccount", trace.WithAttributes(
		attribute.String("account_id", req.GetAccountId()),
	))
	defer span.End()

	if err := validate.Validate(req); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating get account request: %w", err), ierrors.ErrNotValidRequestData))
	}

	serviceAccount, err := c.service.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, c.respondError(ctx, err)
	}

	account := transformer.NewProtoAccountFromAccount(serviceAccount)
	if err := validate.Validate(account); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating get account from account: %w", err), ierrors.ErrNotValidInternalData))
	}

	return &servicev1.GetAccountResponse{Account: account}, nil
}

func (c *Controller) ListAccounts(ctx context.Context, req *servicev1.ListAccountsRequest) (*servicev1.ListAccountsResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.ListAccounts", trace.WithAttributes(
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

	accounts, nextCursor, err := c.service.ListAccounts(ctx, cursor, limit)
	if err != nil {
		return nil, c.respondError(ctx, err)
	}

	listAccountsResponse := &servicev1.ListAccountsResponse{
		NextPageToken: nextCursor,
	}
	for _, account := range accounts {
		listAccountsResponse.Accounts = append(listAccountsResponse.Accounts, transformer.NewProtoAccountFromAccount(account))
	}
	if err := validate.Validate(listAccountsResponse); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating list accounts response from accounts: %w", err), ierrors.ErrNotValidInternalData))
	}

	return listAccountsResponse, nil
}

func (c *Controller) CreateAccount(ctx context.Context, create *servicev1.CreateAccountRequest) (*servicev1.CreateAccountResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.CreateAccount", trace.WithAttributes(
		attribute.String("parent_id", create.GetParentId()),
		attribute.String("name", create.GetName()),
		attribute.String("type", create.GetType().String()),
		attribute.String("basis", create.GetBasis().String()),
	))
	defer span.End()

	if err := validate.Validate(create); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating create account request: %w", err), ierrors.ErrNotValidRequestData))
	}

	serviceAccount := transformer.NewAccountFromProtoCreateAccount(create)

	serviceAccount, err := c.service.CreateAccount(ctx, serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	account := transformer.NewProtoAccountFromAccount(serviceAccount)
	if err := validate.Validate(account); err != nil {
		return nil, c.respondError(ctx, errors.Join(fmt.Errorf("validating account from created account: %w", err), ierrors.ErrNotValidInternalData))
	}
	return &servicev1.CreateAccountResponse{Account: account}, nil
}
