package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hampgoodwin/GoLuca/internal/errors"
	"go.uber.org/zap"
)

type errorResponse struct {
	Description      string `json:"description"`
	ValidationErrors string `json:"validationErrors,omitempty"`
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
		c.respond(w, errorResponse{Description: "bad request data, check request meta data"}, http.StatusBadRequest)
		return
	case errors.HasFlag(err, errors.NotValidRequestData):
		if c.respondOnValidationErrors(w, err, "bad request data.") {
			log.Error("respondError", zap.Error(err))
			return
		}
		c.respond(w, errorResponse{Description: "bad request data."}, http.StatusBadRequest)
		log.Error(
			"incorrect error flag used for case",
			zap.Error(err), zap.String("error_flag", fmt.Sprint(errors.NotValidRequestData)),
		)
		return
	case errors.HasFlag(err, errors.NotFound):
		c.respond(w, nil, http.StatusNotFound)
		return
	case errors.HasFlag(err, errors.NotValidInternalData):
		c.respond(w, errorResponse{Description: "internal data is invalid and failed validation."}, http.StatusInternalServerError)
	case errors.HasFlag(err, errors.NotDeserializable):
		c.respond(w, errorResponse{Description: "provided data passed failed deserialization. If creating a resource, check the request body types."}, http.StatusInternalServerError)
	case errors.HasFlag(err, errors.NotSerializable):
		c.respond(w, errorResponse{Description: "either provided or internal data passed validation, but failed serialization."}, http.StatusInternalServerError)
	case errors.HasFlag(err, errors.NoRelationshipFound):
		c.respond(w, errorResponse{Description: "process which assumed existence of a relationship between data found no relationship. If you are creating data with related data id, those id's do not exist."}, http.StatusBadRequest)
	default:
		log.Error("respondError", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unhandled error response."))
	}
}

func (c *Controller) respondOnValidationErrors(w http.ResponseWriter, err error, description string) bool {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.respond(w, errorResponse{Description: description, ValidationErrors: validationErrors.Error()}, http.StatusBadRequest)
		return true
	}
	return false
}
