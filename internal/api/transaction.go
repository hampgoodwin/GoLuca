package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
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

type transactionsResponse struct {
	Transactions []transaction.Transaction `json:"transactions" validate:"required"`
}

type transactionEntriesResponse struct {
	Entries []transaction.Entry `json:"entries" validate:"required"`
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
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
	transactions, err := service.GetTransactions(ctx, cursorInt, limitInt)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
	res := &transactionsResponse{Transactions: transactions}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transaction_id")
	transactionIDInt, err := strconv.ParseInt(transactionID, 10, 64) // the route regexp handles err cases
	if err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	transaction, err := service.GetTransaction(ctx, transactionIDInt)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
	res := &transactionResponse{Transaction: transaction}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}

func getTransactionEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transaction_id")
	transactionIDInt, err := strconv.ParseInt(transactionID, 10, 64) // the route regexp handles err cases
	if err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	entries, err := service.GetTransactionEntries(ctx, transactionIDInt)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}

	res := transactionEntriesResponse{entries}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &transactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}

	if !req.Transaction.Balanced() {
		err := fmt.Errorf("transaction entires are not balanced\n%s", req.Transaction.Entries)
		encodeError(w, http.StatusBadRequest, err)
		return
	}

	transaction, err := service.CreateTransaction(ctx, req.Transaction)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}

	res := &transactionResponse{Transaction: transaction}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}
