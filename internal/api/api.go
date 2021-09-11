package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"github.com/hampgoodwin/GoLuca/internal/lucalog"
	"go.uber.org/zap"
)

type errorResponse struct {
	Description      string                     `json:"description"`
	ValidationErrors validator.ValidationErrors `json:"validationErrors,omitempty"`
}

func respond(w http.ResponseWriter, i interface{}, statuseCode int) {
	w.WriteHeader(statuseCode)
	if err := json.NewEncoder(w).Encode(i); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("encoding response"))
	}
}

func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.HasFlag(err, errors.NotValidRequest):
		respond(
			w,
			errorResponse{Description: "bad data in request, check request meta data"},
			http.StatusBadRequest,
		)
		return
	case errors.HasFlag(err, errors.NotValidRequestData):
		if respondOnValidationErrors(w, err, "bad request data") {
			return
		}
		respond(w, errorResponse{Description: "bad request data"}, http.StatusBadRequest)
		lucalog.Logger.Error("incorrect error flag used for case",
			zap.String("error", fmt.Sprint(errors.NotValidRequestData)))
		return
	default:
		lucalog.Logger.Error("respondError", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unhandled error response"))
	}
}

func respondOnValidationErrors(w http.ResponseWriter, err error, description string) bool {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		respond(w, errorResponse{Description: description, ValidationErrors: validationErrors}, http.StatusBadRequest)
		return true
	}
	return false
}
