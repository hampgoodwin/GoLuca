package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Accounts   []account.Account `json:"accounts" validated:"required"`
	NextCursor string            `json:"nextCursor,omitempty" validated:"base64"`
}

func (c *Controller) RegisterAccountRoutes(r *chi.Mux) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", c.getAccounts)
		r.Get(fmt.Sprintf("/{accountId:%s}", uuidRegexp), c.getAccount)
		r.Post("/", c.createAccount)
	})
}

func (c *Controller) getAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID := chi.URLParam(r, "accountId")
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
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	if limit == "" {
		limit = "3"
	}
	accounts, nextCursor, err := c.service.GetAccounts(ctx, cursor, limit)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting accounts from service"))
		return
	}
	if err := validate.Validate(accounts); err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "validating accounts", errors.NotValidInternalData))
		return
	}

	res := &accountsResponse{Accounts: accounts, NextCursor: *nextCursor}
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
