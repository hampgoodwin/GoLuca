package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/pkg/transaction"
	"github.com/go-chi/chi"
)

func registerEntryRoute(r *chi.Mux) {
	r.Route("/entries", func(r chi.Router) {
		r.Get("/entries", getEntries) // GET /entries
		// r.Get("/entries/{id:[0-9]+}", getEntry) // GET /entries/1

		// r.Post("/entries", createEntry) // POST /entries
	})
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tx, err := data.DB.BeginTx(ctx, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	txStmt, err := tx.Prepare(`SELECT id, account, change FROM entry;`)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	rows, err := txStmt.QueryContext(ctx)
	var entries []*transaction.Entry
	for rows.Next() {
		entry := &transaction.Entry{}
		rows.Scan(
			&entry.ID,
			&entry.Account,
			&entry.Amount,
		)
		entries = append(entries, entry)
	}
	w.Write([]byte(fmt.Sprintf("%v", entries)))
	return
}
