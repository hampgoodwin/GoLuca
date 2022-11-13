package controller

import (
	"context"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *Controller) GetAccount(ctx context.Context, req *servicev1.GetAccountRequest) (*servicev1.GetAccountResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "grpc.v1.controller.GetAccount", trace.WithAttributes(
		attribute.String("account_id", req.GetAccountId()),
	))
	defer span.End()

	if err := validate.Validate(req); err != nil {
		return nil, c.respondError(ctx, c.log, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating get account request"))
	}

	serviceAccount, err := c.service.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, c.respondError(ctx, c.log, err)
	}

	account := transformer.NewProtoAccountFromAccount(serviceAccount)
	if err := validate.Validate(account); err != nil {
		return nil, c.respondError(ctx, c.log, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating get account from account"))
	}

	return &servicev1.GetAccountResponse{Account: account}, nil
}

func (c *Controller) ListAccounts(ctx context.Context, req *servicev1.ListAccountsRequest) (*servicev1.ListAccountsResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "grpc.v1.controller.ListAccounts", trace.WithAttributes(
		attribute.Int64("page_size", int64(req.GetPageSize())),
		attribute.String("page_token", req.GetPageToken()),
	))
	defer span.End()

	limit, cursor := req.PageSize, req.PageToken
	if limit == 0 {
		limit = 10
	}
	if err := validate.Var(cursor, "omitempty,base64"); err != nil {
		return nil, c.respondError(ctx, c.log, errors.WithErrorMessage(err, errors.NotValidRequest, "invalid cursor or token"))
	}

	accounts, nextCursor, err := c.service.ListAccounts(ctx, cursor, limit)
	if err != nil {
		return nil, c.respondError(ctx, c.log, err)
	}

	listAccountsResponse := &servicev1.ListAccountsResponse{
		NextPageToken: nextCursor,
	}
	for _, account := range accounts {
		listAccountsResponse.Accounts = append(listAccountsResponse.Accounts, transformer.NewProtoAccountFromAccount(account))
	}
	if err := validate.Validate(listAccountsResponse); err != nil {
		return nil, c.respondError(ctx, c.log, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating list accounts response from accounts"))
	}

	return listAccountsResponse, nil
}

func (c *Controller) CreateAccount(ctx context.Context, create *servicev1.CreateAccountRequest) (*servicev1.CreateAccountResponse, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "grpc.v1.controller.CreateAccount", trace.WithAttributes(
		attribute.String("parent_id", create.GetParentId()),
		attribute.String("name", create.GetName()),
		attribute.String("type", create.GetType().String()),
		attribute.String("basis", create.GetBasis().String()),
	))
	defer span.End()

	if err := validate.Validate(create); err != nil {
		return nil, c.respondError(ctx, c.log, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating create account request"))
	}

	serviceAccount := transformer.NewAccountFromProtoCreateAccount(create)

	serviceAccount, err := c.service.CreateAccount(ctx, serviceAccount)
	if err != nil {
		return nil, errors.WithMessage(err, "creating account")
	}

	account := transformer.NewProtoAccountFromAccount(serviceAccount)
	if err := validate.Validate(account); err != nil {
		return nil, c.respondError(ctx, c.log, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating account from created account"))
	}
	return &servicev1.CreateAccountResponse{Account: account}, nil
}
