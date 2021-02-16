package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/pkg/account"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
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
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to start db-transaction for getting an account")
		w.Write([]byte(err.Error()))
		return
	}
	txAccountSelectSmt, err := tx.PrepareContext(ctx,
		`SELECT id, parent_id, name, type, basis
FROM account
WHERE id=$1
;`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to prepare account get statement")
		w.Write([]byte(err.Error()))
		return
	}
	account := account.Account{}
	txAccountSelectSmt.QueryRowContext(ctx, accountID).
		Scan(&account.ID, &account.ParentID, &account.Name, &account.Type, &account.Basis)
	fmt.Println(account)
	accountResp := &accountResponse{Account: &account}
	if err := validator.New().Struct(accountResp); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}
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

// TODO: PAGINATE
func getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get a db-transaction
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to start db-transaction for creating transaction")
		w.Write([]byte(err.Error()))
		return
	}
	txAccountsSelectStmt, err := tx.PrepareContext(ctx, `
SELECT id, parent_id, name, type, basis
FROM account
;`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to prepare account get all statement")
		w.Write([]byte(err.Error()))
		return
	}
	rows, err := txAccountsSelectStmt.QueryContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to get accounts on db query")
		w.Write([]byte(err.Error()))
		return
	}
	accounts := []account.Account{}
	for rows.Next() {
		account := account.Account{}
		rows.Scan(
			&account.ID,
			&account.ParentID,
			&account.Name,
			&account.Type,
			&account.Basis,
		)
		accounts = append(accounts, account)
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
