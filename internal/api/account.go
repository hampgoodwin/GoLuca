package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/account"
	"go.uber.org/zap"
)

type accountRequest struct {
	*account.Account `json:"account" validate:"required"`
}

type accountResponse struct {
	*account.Account `json:"account" validate:"required"`
}

type accountsResponse struct {
	Accounts []account.Account `json:"accounts" validated:"required"`
}

func (c *Controller) RegisterAccountRoutes(r *chi.Mux) {
	r.Get("/accounts", c.getAccounts)
	r.Get("/accounts/{account_id:[0-9]+}", c.getAccount)
	r.Post("/accounts", c.createAccount)
}

func (c *Controller) getAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID := chi.URLParam(r, "account_id")
	c.log.Info("getting account", zap.String("account", accountID))

	account, err := c.service.GetAccount(ctx, accountID)
	if err != nil {
		c.respondError(w, c.log, err)
		return
	}
	if err := validate.Validate(account); err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "validating account", errors.NotValidInternalData))
		return
	}

	res := &accountResponse{Account: account}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get query strings for pagination
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "parsing limit query param", errors.NotValidRequest))
		return
	}

	accounts, err := c.service.GetAccounts(ctx, cursor, limitInt)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting accounts from service"))
		return
	}
	if err := validate.Validate(accounts); err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "validating accounts", errors.NotValidInternalData))
		return
	}

	res := &accountsResponse{Accounts: accounts}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) createAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &accountRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "deserializing request body", errors.NotDeserializable))
		return
	}

	created, err := c.service.CreateAccount(ctx, req.Account)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "creating account in service"))
		return
	}

	res := accountResponse{Account: created}
	w.WriteHeader(http.StatusCreated)
	c.respond(w, res, http.StatusCreated)
}
