package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"
	"github.com/hampgoodwin/errors"
)

func (c *Controller) RegisterTransactionRoutes(r *chi.Mux) {
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", c.listTransactions)
		r.Get(fmt.Sprintf("/{transactionId:%s}", ksuidRegexp), c.getTransaction)
		r.Post("/", c.createTransaction)
	})
}

type transactionRequest struct {
	Transaction httptransaction.CreateTransaction `json:"transaction" validate:"required"`
}

type transactionResponse struct {
	httptransaction.Transaction `json:"transaction" validate:"required"`
}

type transactionsResponse struct {
	Transactions []httptransaction.Transaction `json:"transactions" validate:"required"`
	NextCursor   string                        `json:"nextCursor,omitempty" validate:"base64"`
}

func (c *Controller) listTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	if limit == "" {
		limit = "10"
	}

	transactions, nextCursor, err := c.service.ListTransactions(ctx, cursor, limit)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting transactions from service"))
		return
	}

	responseTransactions := []httptransaction.Transaction{}
	for _, transaction := range transactions {
		responseTransactions = append(responseTransactions, transformer.NewHTTPTransactionFromTransaction(transaction))
	}
	if err := validate.Validate(responseTransactions); err != nil {
		c.respondError(w, c.log, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating http transactions from transaction"))
		return
	}

	res := &transactionsResponse{Transactions: responseTransactions, NextCursor: *nextCursor}
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

	responseTransaction := transformer.NewHTTPTransactionFromTransaction(transaction)
	if err := validate.Validate(responseTransaction); err != nil {
		c.respondError(w, c.log, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating http transaction from transaction"))
		return
	}

	res := &transactionResponse{Transaction: responseTransaction}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) createTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &transactionRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c.respondError(w, c.log, errors.WrapWithErrorMessage(err, errors.NotDeserializable, err.Error(), "deserializing request body"))
		return
	}
	if err := validate.Validate(req); err != nil {
		c.respondError(w, c.log, errors.WithErrorMessage(err, errors.NotValidRequestData, "validating http api transaction request"))
		return
	}

	create, err := transformer.NewTransactionFromHTTPCreateTransaction(req.Transaction)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(errors.WithError(err, errors.NotValidRequest), "transforming transaction from http transaction"))
		return
	}

	created, err := c.service.CreateTransaction(ctx, create)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "creating transaction in service"))
		return
	}

	returning := transformer.NewHTTPTransactionFromTransaction(created)
	if err := validate.Validate(returning); err != nil {
		c.respondError(w, c.log, errors.WithErrorMessage(err, errors.NotValidInternalData, "validating http transaction from transaction"))
		return
	}

	res := &transactionResponse{Transaction: returning}
	c.respond(w, res, http.StatusCreated)
}
