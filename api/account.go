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
	r.Get("/accounts/{id:[0-9]+}", getAccount)
	r.Post("/accounts", createAccount)
}

func getAccount(w http.ResponseWriter, r *http.Request) {

}

// TODO: PAGINATE
func getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get a db-transaction
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to start db-transaction for creating transaction"))
		return
	}
	txAccountSelectStmt, err := tx.PrepareContext(ctx, `
SELECT id, parent_id, name, type, basis
FROM account
;`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to prepare account get all statement"))
		return
	}
	rows, err := txAccountSelectStmt.QueryContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get accounts on db query"))
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
		w.Write([]byte("failed to encode entries response"))
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
		w.Write([]byte(fmt.Sprintf("failed to decode body\n%s\n", err.Error())))
		return
	}
	if err := validator.New().Struct(aReq); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(string(err.Error())))
			return
		}
	}
	// get a db-transaction
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to start db-transaction for creating transaction"))
		return
	}
	txInsertAccountStmt, err := tx.PrepareContext(ctx, `INSERT INTO account(parent_id, name, type, basis)
	VALUES($1, $2, $3, $4)
	RETURNING id, parent_id, name, type, basis;`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to prepare account insert statement")
		w.Write([]byte(err.Error()))
		tx.Rollback()
		return
	}
	account := &account.Account{}
	txInsertAccountStmt.QueryRowContext(ctx, aReq.ParentID, aReq.Name, aReq.Type, aReq.Basis).Scan(
		&account.ID,
		&account.ParentID,
		&account.Name,
		&account.Type,
		&account.Basis,
	)
	aRes := &accountResponse{Account: account}
	if err := validator.New().Struct(aRes); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			err = errors.Wrap(err, "response type failed validation")
			w.Write([]byte(err.Error()))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(&aRes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode createAccount response"))
		tx.Rollback()
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusCreated)
	return
}
