package api

import (
	"encoding/json"
	"net/http"
)

func encode(w http.ResponseWriter, i interface{}) error {
	if err := json.NewEncoder(w).Encode(i); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
