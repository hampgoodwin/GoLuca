package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/lucalog"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
)

func registerEntryRoute(r *chi.Mux) {
	r.Get("/entries", getEntries) // GET /entries
}

type entriesResponse struct {
	Entries []transaction.Entry `json:"entries,omitempty"`
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	lucalog.Logger.Info("getting entries")
	ctx := r.Context()

	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		respondError(w, errors.WrapFlag(err, "parsing limit query param", errors.NotValidRequest))
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		respondError(w, errors.WrapFlag(err, "parsing cursor query param", errors.NotValidRequest))
		return
	}
	entries, err := service.GetEntries(ctx, cursorInt, limitInt)
	if err != nil {
		respondError(w, errors.Wrap(err, "getting entries from service"))
	}

	res := &entriesResponse{Entries: entries}
	respond(w, res, http.StatusOK)
}
