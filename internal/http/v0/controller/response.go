package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hampgoodwin/errors"
	"go.uber.org/zap"
)

type ErrorResponse struct {
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
	log.Error("responding", zap.Error(err))

	var statuscode int
	var message string
	var msg errors.Message
	if errors.As(err, &msg) {
		message = msg.Value
	}

	switch {
	case errors.Is(err, errors.NotKnown):
		statuscode = http.StatusInternalServerError
	case errors.Is(err, errors.NotValidRequest):
		statuscode = http.StatusBadRequest
		message = "bad request data, check request meta data."
	case errors.Is(err, errors.NotValidRequestData):
		if message == "" {
			message = "bad request data"
		}
		if c.respondOnValidationErrors(w, err, message) {
			return
		}
		c.respond(w, ErrorResponse{Description: message}, http.StatusBadRequest)
		log.Error(
			"incorrect error flag used for case",
			zap.Error(err), zap.String("error_flag", fmt.Sprint(errors.NotValidRequestData)),
		)
		return
	case errors.Is(err, errors.NotFound):
		statuscode = http.StatusNotFound
	case errors.Is(err, errors.NotValidInternalData):
		statuscode = http.StatusInternalServerError
		message = "internal data is invalid and failed validation."
	case errors.Is(err, errors.NotDeserializable):
		statuscode = http.StatusInternalServerError
		message = "provided data passed failed deserialization. If creating a resource, check the request body types."
	case errors.Is(err, errors.NotSerializable):
		statuscode = http.StatusInternalServerError
		message = "either provided or internal data passed validation, but failed serialization."
	case errors.Is(err, errors.NoRelationshipFound):
		statuscode = http.StatusBadRequest
		message = "process which assumed existence of a relationship between data found no relationship. If you are creating data with related data id, those id's do not exist."
	default:
		log.Error("respondError", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unhandled error response."))
	}
	if errors.As(err, &msg) {
		if msg.Value != "" {
			message = msg.Value
		}
	}
	c.respond(w, ErrorResponse{Description: message}, statuscode)
}

func (c *Controller) respondOnValidationErrors(w http.ResponseWriter, err error, message string) bool {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.respond(w, ErrorResponse{Description: message, ValidationErrors: validationErrors.Error()}, http.StatusBadRequest)
		return true
	}
	return false
}
