package controller

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (c *Controller) RegisterEntryRoutes(r *chi.Mux) {
	r.Get("/entries", c.getEntries)
}

type entriesResponse struct {
	Entries    []transaction.Entry `json:"entries" validated:"required"`
	NextCursor string              `json:"nextCursor,omitempty" validated:"base64"`
}

func (c *Controller) getEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	if limit == "" {
		limit = "10"
	}

	entries, nextCursor, err := c.service.GetEntries(ctx, cursor, limit)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting entries from service"))
		return
	}

	res := &entriesResponse{Entries: entries, NextCursor: *nextCursor}
	c.respond(w, res, http.StatusOK)
}
