package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"go.uber.org/zap"
)

type errorResponse struct {
	Description      string                     `json:"description"`
	ValidationErrors validator.ValidationErrors `json:"validationErrors,omitempty"`
}

func (c *Controller) respond(w http.ResponseWriter, i interface{}, statuseCode int) {
	w.WriteHeader(statuseCode)
	if err := json.NewEncoder(w).Encode(i); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("encoding response"))
	}
}

func (c *Controller) respondError(w http.ResponseWriter, log *zap.Logger, err error) {
	switch {
	case errors.HasFlag(err, errors.NotValidRequest):
		c.respond(w, errorResponse{Description: "bad data in request, check request meta data"}, http.StatusBadRequest)
		return
	case errors.HasFlag(err, errors.NotValidRequestData):
		if c.respondOnValidationErrors(w, err, "bad request data") {
			return
		}
		c.respond(w, errorResponse{Description: "bad request data"}, http.StatusBadRequest)
		log.Error("incorrect error flag used for case",
			zap.String("error", fmt.Sprint(errors.NotValidRequestData)))
		return
	case errors.HasFlag(err, errors.NotFound):
		c.respond(w, errorResponse{}, http.StatusNotFound)
		return
	case errors.HasFlag(err, errors.NotValidInternalData):
		c.respond(w, errorResponse{Description: "internal data is invalid and failed validation"}, http.StatusInternalServerError)
	case errors.HasFlag(err, errors.NotDeserializable):
		c.respond(w, errorResponse{Description: "provided data passed validation but failed deserialization to internal type"}, http.StatusInternalServerError)
	case errors.HasFlag(err, errors.NotSerializable):
		c.respond(w, errorResponse{Description: "either provided or internal data passed validation, but failed serialization"}, http.StatusInternalServerError)
	default:
		log.Error("respondError", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unhandled error response"))
	}
}

func (c *Controller) respondOnValidationErrors(w http.ResponseWriter, err error, description string) bool {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.respond(w, errorResponse{Description: description, ValidationErrors: validationErrors}, http.StatusBadRequest)
		return true
	}
	return false
}
