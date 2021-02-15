package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/pkg/transaction"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func registerTransactionRoute(r *chi.Mux) {
	r.Post("/transactions", createTransaction)
}

type transactionRequest struct {
	*transaction.Transaction `json:"transaction" validate:"required"`
}

type transactionResponse struct {
	*transaction.Transaction `json:"transaction" validate:"required"`
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tReq := &transactionRequest{}

	if err := json.NewDecoder(r.Body).Decode(tReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(err, "failed to decode body")
		w.Write([]byte(err.Error()))
		return
	}

	// get a db-transaction
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to start db-transaction for creating transaction")
		w.Write([]byte(err.Error()))
		return
	}

	// insert the transaction
	txInsertTransactionStmt, err := tx.PrepareContext(ctx, `
INSERT INTO transaction (description) VALUES ($1)
RETURNING id, description;`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to prepare transaction insert statement")
		w.Write([]byte(err.Error()))
		return
	}
	transactionCreated := transaction.Transaction{}
	txInsertTransactionStmt.QueryRowContext(ctx, tReq.Transaction.Description).
		Scan(&transactionCreated.ID, transactionCreated.Description)

	// insert the entries
	for _, entry := range tReq.Transaction.Entries {
		txInsertEntryStmt, err := tx.PrepareContext(ctx, `
INSERT INTO entry(transaction_id, account_id, amount) VALUES ($1, $2, $3)
RETURNING id, transaction_id, account_id, amount;`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = errors.Wrap(err, "failed to prepare entry insert statement")
			w.Write([]byte(err.Error()))
			return
		}
		entryCreated := transaction.Entry{}
		txInsertEntryStmt.QueryRowContext(ctx, transactionCreated.ID, entry.AccountID, entry.Amount).
			Scan(&entryCreated.ID, &entryCreated.TransactionID, &entryCreated.AccountID, &entryCreated.Amount)
		fmt.Println(entryCreated)
		transactionCreated.Entries = append(transactionCreated.Entries, entryCreated)
	}
	tx.Commit()

	tRes := &transactionResponse{Transaction: &transactionCreated}
	if err := json.NewEncoder(w).Encode(tRes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to encode response")
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusCreated)
	return
}
