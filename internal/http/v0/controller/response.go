package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"

	"github.com/go-playground/validator/v10"
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
	defer span.SetStatus(otelcodes.Error, message)

	var notFoundErr ierrors.NotFoundErr
	if errors.As(err, &notFoundErr) {
		c.respond(w, ErrorResponse{Description: notFoundErr.Error()}, http.StatusNotFound)
		return
	}
	var notValidCursorErr ierrors.NotValidCursorErr
	if errors.As(err, &notValidCursorErr) {
		c.respond(w, ErrorResponse{Description: notValidCursorErr.Error()}, http.StatusBadRequest)
		return
	}

	switch {
	case errors.Is(err, ierrors.ErrNotKnown):
		statuscode = http.StatusInternalServerError
	case errors.Is(err, ierrors.ErrNotValidRequest):
		statuscode = http.StatusBadRequest
		message = "bad request data, check request meta data."
	case errors.Is(err, ierrors.ErrNotValidRequestData):
		if message == "" {
			message = "bad request data"
		}
		if c.respondOnValidationErrors(w, err, message) {
			return
		}
		c.respond(w, ErrorResponse{Description: message}, http.StatusBadRequest)
		return
	case errors.Is(err, ierrors.ErrNotFound):
		statuscode = http.StatusNotFound
		message = "not found"
	case errors.Is(err, ierrors.ErrNotValidInternalData):
		statuscode = http.StatusInternalServerError
		message = "internal data is invalid and failed validation."
	case errors.Is(err, ierrors.ErrNotDeserializable):
		statuscode = http.StatusInternalServerError
		message = "provided data passed failed deserialization. If creating a resource, check the request body types."
	case errors.Is(err, ierrors.ErrNotSerializable):
		statuscode = http.StatusInternalServerError
		message = "either provided or internal data passed validation, but failed serialization."
	case errors.Is(err, ierrors.ErrNoRelationshipFound):
		statuscode = http.StatusBadRequest
		message = "process which assumed existence of a relationship between data found no relationship. If you are creating data with related data id, those id's do not exist."
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("unhandled error response."))
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
