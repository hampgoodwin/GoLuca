package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (c *Controller) RegisterTransactionRoutes(r *chi.Mux) {
	r.Get("/transactions", c.getTransactions)
	r.Get("/transactions/{transaction_id:[0-9]+}", c.getTransaction)
	r.Get("/transactions/{transaction_id:[0-9]+}/entries", c.getTransactionEntries)
	r.Post("/transactions", c.createTransaction)
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

func (c *Controller) getTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get query strings for pagination
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "parsing limit query parameter", errors.NotValidRequest))
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "parsing cursor query parameter", errors.NotValidRequest))
		return
	}
	transactions, err := c.service.GetTransactions(ctx, cursorInt, limitInt)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting transactions from service"))
		return
	}
	res := &transactionsResponse{Transactions: transactions}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) getTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transaction_id")
	transactionIDInt, err := strconv.ParseInt(transactionID, 10, 64) // the route regexp handles err cases
	if err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "parsing transactionID URL parameter", errors.NotValidRequest))
		return
	}
	transaction, err := c.service.GetTransaction(ctx, transactionIDInt)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting transaction from service"))
		return
	}
	res := &transactionResponse{Transaction: transaction}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) getTransactionEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transaction_id")
	transactionIDInt, err := strconv.ParseInt(transactionID, 10, 64) // the route regexp handles err cases
	if err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "parsing transactionID URL parameter", errors.NotValidRequest))
		return
	}
	entries, err := c.service.GetTransactionEntries(ctx, transactionIDInt)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting transaction entries from service"))
		return
	}

	res := transactionEntriesResponse{entries}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) createTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &transactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "deserializing request body", errors.NotValidRequestData))
		return
	}

	if !req.Transaction.Balanced() {
		err := fmt.Errorf("transaction entires are not balanced\n%s", req.Transaction.Entries)
		c.respondError(w, c.log, errors.WrapFlag(err, "validating transaction request is balanced", errors.NotValidRequest))
		return
	}

	transaction, err := c.service.CreateTransaction(ctx, req.Transaction)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "creating transaction in service"))
		return
	}

	res := &transactionResponse{Transaction: transaction}
	c.respond(w, res, http.StatusCreated)
}
