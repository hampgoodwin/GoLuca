package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/data"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/lucalog"
	"github.com/hampgoodwin/GoLuca/internal/service"
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

func registerAccountRoutes(r *chi.Mux) {
	r.Get("/accounts", getAccounts)
	r.Get("/accounts/{account_id:[0-9]+}", getAccount)
	r.Post("/accounts", createAccount)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID := chi.URLParam(r, "account_id")
	lucalog.Logger.Info("getting account", zap.String("account", accountID))
	accountIDInt, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		err = errors.WrapFlag(err, "parsing account id query parameter", errors.NotValidRequest)
		respondError(w, err)
	}

	account, err := service.GetAccount(ctx, accountIDInt)
	if err != nil {
		respondError(w, err)
		return
	}
	if err := data.Validate(account); err != nil {
		respondError(w, errors.WrapFlag(err, "validating account", errors.NotValidInternalData))
		return
	}

	res := &accountResponse{Account: account}
	respond(w, res, http.StatusOK)
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get query strings for pagination
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		respondError(w, errors.WrapFlag(err, "parsing limit query param", errors.NotValidRequest))
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		respondError(w, errors.WrapFlag(err, "parsing cursor int query param", errors.NotValidRequest))
		return
	}

	accounts, err := service.GetAccounts(ctx, cursorInt, limitInt)
	if err != nil {
		respondError(w, errors.Wrap(err, "getting accounts from service"))
		return
	}
	if err := data.Validate(accounts); err != nil {
		respondError(w, errors.WrapFlag(err, "validating accounts", errors.NotValidInternalData))
		return
	}

	res := &accountsResponse{Accounts: accounts}
	respond(w, res, http.StatusOK)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &accountRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respondError(w, errors.WrapFlag(err, "deserializing request body", errors.NotDeserializable))
		return
	}

	created, err := service.CreateAccount(ctx, req.Account)
	if err != nil {
		respondError(w, errors.Wrap(err, "creating account in service"))
		return
	}

	res := accountResponse{Account: created}
	w.WriteHeader(http.StatusCreated)
	respond(w, res, http.StatusCreated)
}
