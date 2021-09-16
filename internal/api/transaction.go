package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (c *Controller) RegisterTransactionRoutes(r *chi.Mux) {
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", c.getTransactions)
		r.Get(fmt.Sprintf("/{transactionId:%s}", uuidRegexp), c.getTransaction)
		r.Post("/", c.createTransactionAndEntries)
	})
}

type Transaction struct {
	Description string    `json:"description" validate:"required"`
	Entries     []Entry   `json:"entries,omitempty" validate:"dive,gte=1"`
	CreatedAt   time.Time `json:"createdAt" validate:"required"`
}

type Entry struct {
	Description   string    `json:"description"`
	DebitAccount  string    `json:"debitAccount" validate:"required,uuid4"`
	CreditAccount string    `json:"creditAccount" validate:"required,uuid4"`
	Amount        Amount    `json:"amount" validate:"required"`
	CreatedAt     time.Time `json:"createdAt" validate:"required"`
}

type Amount struct {
	Value    string `json:"value" validate:"gte=0"`
	Currency string `json:"currency" validate:"len=3,alpha"`
}

type transactionRequest struct {
	Transaction `json:"transaction" validate:"required"`
}

type transactionResponse struct {
	*transaction.Transaction `json:"transaction" validate:"required"`
}

type transactionsResponse struct {
	Transactions []transaction.Transaction `json:"transactions" validate:"required"`
	NextCursor   string                    `json:"nextCursor,omitempty" validate:"base64"`
}

func (c *Controller) getTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	if limit == "" {
		limit = "3"
	}
	transactions, nextCursor, err := c.service.GetTransactions(ctx, cursor, limit)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting transactions from service"))
		return
	}
	res := &transactionsResponse{Transactions: transactions, NextCursor: *nextCursor}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) getTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionID := chi.URLParam(r, "transactionId")
	transaction, err := c.service.GetTransaction(ctx, transactionID)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting transaction from service"))
		return
	}
	res := &transactionResponse{Transaction: transaction}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) createTransactionAndEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &transactionRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "deserializing request body", errors.NotValidRequestData))
		return
	}

	transaction, err := c.service.CreateTransactionAndEntries(ctx, req.Transaction)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "creating transaction in service"))
		return
	}

	res := &transactionResponse{Transaction: transaction}
	c.respond(w, res, http.StatusCreated)
}
