package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/pkg/account"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
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
	account, err := data.GetAccount(ctx, accIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to query account from database")
		w.Write([]byte(err.Error()))
		return
	}

	accountResp := &accountResponse{Account: account}
	if err := data.Validate(accountResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "data validation for gathered type failed")
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(accountResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to encode account response")
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get query strings for pagination
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(err, "required query string limit must be integer")
		w.Write([]byte(err.Error()))
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(err, "required query string cursor must be integer")
		w.Write([]byte(err.Error()))
		return
	}

	accounts, err := data.GetAccounts(ctx, limitInt, cursorInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrapf(err, "failed to get accounts from database with limit %d, offset %d", limitInt, cursorInt)
		w.Write([]byte(err.Error()))
		return
	}

	accountsResp := &accountsResponse{Accounts: accounts}
	if err := json.NewEncoder(w).Encode(accountsResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to encode entries response")
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	aReq := &accountRequest{}
	if err := json.NewDecoder(r.Body).Decode(aReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(err, "failed to decode body")
		w.Write([]byte(err.Error()))
		return
	}
	if err := data.Validate(aReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	acc, err := data.CreateAccount(ctx, aReq.Account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	aRes := accountResponse{Account: acc}
	if err := encode(w, aRes); err != nil {
		err = errors.Wrap(err, "failed to encode createAccount response")
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
