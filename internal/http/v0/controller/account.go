package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	httpaccount "github.com/hampgoodwin/GoLuca/pkg/http/v0/account"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/hampgoodwin/errors"
)

type accountRequest struct {
	Account httpaccount.CreateAccount `json:"account" validate:"required"`
}

type accountResponse struct {
	httpaccount.Account `json:"account" validate:"required"`
}

type accountsResponse struct {
	Accounts   []httpaccount.Account `json:"accounts" validate:"required"`
	NextCursor string                `json:"nextCursor,omitempty" validate:"cursor"`
}

func (c *Controller) RegisterAccountRoutes(r *chi.Mux) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", c.listAccounts)
		r.Get(fmt.Sprintf("/{accountId:%s}", ksuidRegexp), c.getAccount)
		r.Post("/", c.createAccount)
	})
}

func (c *Controller) getAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.getAccount")
	defer span.End()

	accountID := chi.URLParam(r, "accountId")
	span.SetAttributes(attribute.String("account_id", accountID))

	account, err := c.service.GetAccount(ctx, accountID)
	if err != nil {
		c.respondError(ctx, w, err)
		return
	}

	responseAccount := transformer.NewHTTPAccountFromAccount(account)
	if err := validate.Validate(responseAccount); err != nil {
		c.respondError(ctx, w, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating http account from account"))
		return
	}

	res := &accountResponse{Account: responseAccount}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) listAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.listAccounts")
	defer span.End()

	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	if limit == "" {
		limit = "10"
	}
	limitUInt64, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		c.respondError(ctx, w, errors.Wrap(err, "converting page size"))
	}
	if cursor == "\"\"" {
		cursor = ""
	}
	span.SetAttributes(
		attribute.String("cursor", cursor),
		attribute.Int64("limit", int64(limitUInt64)),
	)
	if err := validate.Var(cursor, "omitempty,base64"); err != nil {
		c.respondError(ctx, w, errors.WithErrorMessage(err, errors.NotValidRequest, "invalid cursor or token"))
		return
	}

	accounts, nextCursor, err := c.service.ListAccounts(ctx, cursor, limitUInt64)
	if err != nil {
		c.respondError(ctx, w, errors.Wrap(err, "getting accounts from service"))
		return
	}

	responseAccounts := []httpaccount.Account{}
	for _, account := range accounts {
		responseAccounts = append(responseAccounts, transformer.NewHTTPAccountFromAccount(account))
	}
	if err := validate.Validate(responseAccounts); err != nil {
		c.respondError(ctx, w, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating http accounts from accounts"))
		return
	}

	res := &accountsResponse{Accounts: responseAccounts, NextCursor: nextCursor}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) createAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.createAccount")
	defer span.End()
	req := &accountRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c.respondError(ctx, w, errors.WrapWithErrorMessage(err, errors.NotDeserializable, err.Error(), "deserializing request body"))
		return
	}
	if err := validate.Validate(req); err != nil {
		c.respondError(ctx, w, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating http api create account request"))
		return
	}

	create := transformer.NewAccountFromHTTPCreateAccount(req.Account)

	created, err := c.service.CreateAccount(ctx, create)
	if err != nil {
		c.respondError(ctx, w, errors.Wrap(err, "creating account in service"))
		return
	}

	responseCreated := transformer.NewHTTPAccountFromAccount(created)
	if err := validate.Validate(responseCreated); err != nil {
		c.respondError(ctx, w, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating http account from account"))
	}

	res := accountResponse{Account: responseCreated}
	c.respond(w, res, http.StatusCreated)
}
