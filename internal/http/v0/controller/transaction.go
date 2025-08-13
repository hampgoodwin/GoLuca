package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/GoLuca/internal/transformer"
	"github.com/hampgoodwin/GoLuca/internal/validate"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"
	httptransaction "github.com/hampgoodwin/GoLuca/pkg/http/v0/transaction"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (c *Controller) RegisterTransactionRoutes(r *chi.Mux) {
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", c.listTransactions)
		r.Get(fmt.Sprintf("/{transactionId:%s}", uuid7Regexp), c.getTransaction)
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
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.listTransactions")
	defer span.End()
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	if limit == "" {
		limit = "10"
	}
	limitUInt64, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		c.respondError(ctx, w, fmt.Errorf("converting page size: %w", err))
	}

	transactions, nextCursor, err := c.service.ListTransactions(ctx, cursor, limitUInt64)
	if err != nil {
		c.respondError(ctx, w, fmt.Errorf("getting transactions from service: %w", err))
		return
	}

	responseTransactions := []httptransaction.Transaction{}
	for _, transaction := range transactions {
		responseTransactions = append(responseTransactions, transformer.NewHTTPTransactionFromTransaction(transaction))
	}
	if err := validate.Validate(responseTransactions); err != nil {
		c.respondError(ctx, w, errors.Join(fmt.Errorf("validating http transactions from transaction: %w", err), ierrors.ErrNotValidInternalData))
		return
	}

	res := &transactionsResponse{Transactions: responseTransactions, NextCursor: nextCursor}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) getTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.getTransaction")
	defer span.End()
	transactionID := chi.URLParam(r, "transactionId")
	span.SetAttributes(attribute.String("transaction_id", transactionID))

	transaction, err := c.service.GetTransaction(ctx, transactionID)
	if err != nil {
		c.respondError(ctx, w, fmt.Errorf("getting transaction from service: %w", err))
		return
	}

	responseTransaction := transformer.NewHTTPTransactionFromTransaction(transaction)
	if err := validate.Validate(responseTransaction); err != nil {
		c.respondError(ctx, w, errors.Join(fmt.Errorf("validating http transaction from transaction: %w", err), ierrors.ErrNotValidInternalData))
		return
	}

	res := &transactionResponse{Transaction: responseTransaction}
	c.respond(w, res, http.StatusOK)
}

func (c *Controller) createTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.createTransaction")
	defer span.End()
	req := &transactionRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c.respondError(ctx, w, errors.Join(fmt.Errorf("deserializing request body: %w", err), ierrors.ErrNotDeserializable))
		return
	}
	if err := validate.Validate(req); err != nil {
		c.respondError(ctx, w, errors.Join(fmt.Errorf("validating http api transaction request: %w", err), ierrors.ErrNotValidRequestData))
		return
	}

	create, err := transformer.NewTransactionFromHTTPCreateTransaction(req.Transaction)
	if err != nil {
		c.respondError(ctx, w, fmt.Errorf("transforming transaction from http transaction: %w", errors.Join(err, ierrors.ErrNotValidRequest)))
		return
	}

	created, err := c.service.CreateTransaction(ctx, create)
	if err != nil {
		c.respondError(ctx, w, fmt.Errorf("creating transaction in service: %w", err))
		return
	}

	returning := transformer.NewHTTPTransactionFromTransaction(created)
	if err := validate.Validate(returning); err != nil {
		c.respondError(ctx, w, errors.Join(fmt.Errorf("validating http transaction from transaction: %w", err), ierrors.ErrNotValidInternalData))
		return
	}

	res := &transactionResponse{Transaction: returning}
	c.respond(w, res, http.StatusCreated)
}
