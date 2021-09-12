package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func (c *Controller) RegisterEntryRoutes(r *chi.Mux) {
	r.Get("/entries", c.getEntries) // GET /entries
}

type entriesResponse struct {
	Entries []transaction.Entry `json:"entries" validated:"required"`
}

func (c *Controller) getEntries(w http.ResponseWriter, r *http.Request) {
	c.log.Info("getting entries")
	ctx := r.Context()

	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		c.respondError(w, c.log, errors.WrapFlag(err, "parsing limit query param", errors.NotValidRequest))
		return
	}

	entries, err := c.service.GetEntries(ctx, cursor, limitInt)
	if err != nil {
		c.respondError(w, c.log, errors.Wrap(err, "getting entries from service"))
	}

	res := &entriesResponse{Entries: entries}
	c.respond(w, res, http.StatusOK)
}
