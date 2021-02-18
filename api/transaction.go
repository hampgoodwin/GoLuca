package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/pkg/transaction"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func registerTransactionRoute(r *chi.Mux) {
	r.Get("/transactions", getTransactions)
	r.Get("/transactions/{transaction_id:[0-9]+}", getTransaction)
	r.Get("/transactions/{transaction_id:[0-9]+}/entries", getTransactionEntries)
	r.Post("/transactions", createTransaction)
}

type transactionRequest struct {
	*transaction.Transaction `json:"transaction" validate:"required"`
}

type transactionResponse struct {
	*transaction.Transaction `json:"transaction" validate:"required"`
}

type transactionEntriesResponse struct {
	Entries []transaction.Entry `json:"entries" validate:"required"`
}

func getTransactions(w http.ResponseWriter, r *http.Request) {

}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transaction_id")
	transactionIDInt, _ := strconv.ParseInt(transactionID, 10, 64) // the route regexp handles err cases
	transaction, err := data.GetTransaction(ctx, transactionIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrapf(err, "failed to get transaction by id %d from database", transactionIDInt)
		w.Write([]byte(err.Error()))
		return
	}
	transactionResp := &transactionResponse{Transaction: transaction}
	if err := json.NewEncoder(w).Encode(transactionResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to encode response body")
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getTransactionEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transaction_id")
	transactionIDInt, _ := strconv.ParseInt(transactionID, 10, 64) // the route regexp handles err cases
	entries, err := data.GetEntriesByTransactionID(ctx, transactionIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := errors.Wrap(err, "failed to get entries from db")
		w.Write([]byte(err.Error()))
		return
	}

	tRes := transactionEntriesResponse{entries}
	if err := json.NewEncoder(w).Encode(tRes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to encode response")
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
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

	if !tReq.Transaction.Balanced() {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("transaction entires are not balanced\n%s", tReq.Transaction.Entries)
		w.Write([]byte(err.Error()))
		return
	}

	trans, err := data.CreateTransaction(ctx, tReq.Transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to create transaction in db")
		w.Write([]byte(err.Error()))
		return
	}

	tRes := &transactionResponse{Transaction: trans}
	if err := json.NewEncoder(w).Encode(tRes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrap(err, "failed to encode response")
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
