package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/pkg/transaction"
	"github.com/go-chi/chi"
)

func registerEntryRoute(r *chi.Mux) {
	r.Get("/entries", getEntries) // GET /entries
}

type entriesResponse struct {
	Entries transaction.Entries `json:"entries,omitempty"`
}

// TODO: PAGINATE
func getEntries(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	txStmt, err := tx.Prepare(`SELECT id, account, amount FROM entry;`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get entries on db query"))
		return
	}
	rows, err := txStmt.QueryContext(ctx)
	var entries []transaction.Entry
	for rows.Next() {
		entry := transaction.Entry{}
		rows.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.Amount,
		)
		entries = append(entries, entry)
	}
	entriesResp := &entriesResponse{Entries: entries}
	if err := json.NewEncoder(w).Encode(entriesResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode entries response"))
		return
	}
	return
}
