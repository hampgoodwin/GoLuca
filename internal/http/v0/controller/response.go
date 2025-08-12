package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/errors"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
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

func (c *Controller) respondError(ctx context.Context, w http.ResponseWriter, err error) {
	_, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.http.v0.controller.respondError")
	defer span.End()
	span.RecordError(err)

	var statuscode int
	var message string
	var msg errors.Message
	if errors.As(err, &msg) {
		message = msg.Value
	}
	span.SetStatus(otelcodes.Error, message)

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
