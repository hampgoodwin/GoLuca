package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/hampgoodwin/GoLuca/internal/lucalog"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/hampgoodwin/GoLuca/pkg/transaction"
	"go.uber.org/zap"
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
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	cursorInt, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		encodeError(w, http.StatusBadRequest, err)
		return
	}
	entries, err := service.GetEntries(ctx, cursorInt, limitInt)
	if err != nil {
		encodeError(w, http.StatusInternalServerError, err)
	}
	lucalog.Logger.Info("got entries", zap.Int("count", len(entries)))
	res := &entriesResponse{Entries: entries}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}
