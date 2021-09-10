package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"github.com/pkg/errors"
)

func registerEntryRoute(r *chi.Mux) {
	r.Get("/entries", getEntries) // GET /entries
}

type entriesResponse struct {
	Entries []transaction.Entry `json:"entries,omitempty"`
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, cursor := r.URL.Query().Get("limit"), r.URL.Query().Get("cursor")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(err, "required query string limit must be integer")
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(err, "required query string cursor must be integer")
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	entries, err := service.GetEntries(ctx, cursorInt, limitInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrapf(err, "failed to get entries with limit %d, offset %d", limitInt, cursorInt)
		_, _ = w.Write([]byte(err.Error()))
	}
	entriesResp := &entriesResponse{Entries: entries}
	if err := json.NewEncoder(w).Encode(entriesResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to encode entries response"))
		return
	}
}
