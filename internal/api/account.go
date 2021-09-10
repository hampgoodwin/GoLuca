package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/data"
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
	accIDInt, _ := strconv.ParseInt(accountID, 10, 64) // we ignore the err bc the route regexp filters already
	account, err := service.GetAccount(ctx, accIDInt)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}

	res := &accountResponse{Account: account}
	if err := data.Validate(res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := encode(w, res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get query strings for pagination
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}

	accounts, err := service.GetAccounts(ctx, cursorInt, limitInt)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}

	res := &accountsResponse{Accounts: accounts}
	if err := encode(w, res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &accountRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	if err := data.Validate(req); err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	acc, err := service.CreateAccount(ctx, req.Account)
	if err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	res := accountResponse{Account: acc}
	w.WriteHeader(http.StatusCreated)
	if err := encode(w, res); err != nil {
		// TODO split "encodoing" and "responding/writing" code
		encodeError(w, http.StatusInternalServerError, err)
		lucalog.Logger.Error("encoding and writing response", zap.Error(err))
		return
	}
}
