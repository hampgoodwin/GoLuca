package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func encode(w http.ResponseWriter, i interface{}) error {
	if err := json.NewEncoder(w).Encode(i); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func encodeError(w http.ResponseWriter, httpStatusCode int, err error) {
	body := errorResponse{Error: err.Error()}
	w.WriteHeader(httpStatusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error encoding an error response: %s", err.Error())))
	}
}
