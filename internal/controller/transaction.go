package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/httpapi"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (c *Controller) RegisterTransactionRoutes(r *chi.Mux) {
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", c.getTransactions)
		r.Get(fmt.Sprintf("/{transactionId:%s}", uuidRegexp), c.getTransaction)
		r.Post("/", c.createTransactionAndEntries)
	})
}

type transactionRequest struct {
	httpapi.Transaction `json:"transaction" validate:"required"`
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
		limit = "10"
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
		c.respondError(w, c.log, errors.FlagWrap(err, errors.NotDeserializable, err.Error(), "deserializing request body"))
		return
	}

	if err := validate.Validate(req); err != nil {
		c.respondError(w, c.log, errors.FlagWrap(err, errors.NotValidRequestData, "", "validating http api transaction request"))
		return
	}

	create, err := transformer.NewTransactionFromHTTPTransaction(req.Transaction)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "transforming http api transaction to transaction"))
		return
	}

	rtrn, err := c.service.CreateTransactionAndEntries(ctx, &create)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "creating transaction in service"))
		return
	}

	res := &transactionResponse{Transaction: rtrn}
	c.respond(w, res, http.StatusCreated)
}
