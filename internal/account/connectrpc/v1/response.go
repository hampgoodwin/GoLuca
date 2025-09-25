package connect

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"
)

func (h *Handler) respondError(ctx context.Context, err error) error {
	_, span := otel.Tracer(meta.ServiceName).Start(ctx, "internal.grpc.v1.controller.respondError")
	defer span.End()
	span.RecordError(err)

	var statuscode codes.Code
	var message string
	defer span.SetStatus(otelcodes.Error, message)

	var notFoundErr ierrors.NotFoundErr
	if errors.As(err, &notFoundErr) {
		statuscode = codes.NotFound
		message = notFoundErr.Error()
		return status.Error(statuscode, message)
	}
	var notValidTokenErr ierrors.NotValidTokenErr
	if errors.As(err, &notValidTokenErr) {
		statuscode = codes.InvalidArgument
		message = notValidTokenErr.Error()
		return status.Error(statuscode, message)
	}

	switch {
	case errors.Is(err, ierrors.ErrNotKnown):
		statuscode = codes.Unknown
	case errors.Is(err, ierrors.ErrNotValidRequest):
		statuscode = codes.InvalidArgument
		message = "bad request data, check request meta data."
	case errors.Is(err, ierrors.ErrNotValidRequestData):
		statuscode = codes.InvalidArgument
		if message == "" {
			message = "bad request data"
		}
	case errors.Is(err, ierrors.ErrNotFound):
		statuscode = codes.NotFound
	case errors.Is(err, ierrors.ErrNotValidInternalData):
		statuscode = codes.Internal
		message = "internal data is invalid and failed validation."
	case errors.Is(err, ierrors.ErrNotDeserializable):
		statuscode = codes.Internal
		message = "provided data passed failed deserialization. If creating a resource, check the request body types."
	case errors.Is(err, ierrors.ErrNotSerializable):
		statuscode = codes.Internal
		message = "either provided or internal data passed validation, but failed serialization."
	case errors.Is(err, ierrors.ErrNoRelationshipFound):
		statuscode = codes.InvalidArgument
		message = "process which assumed existence of a relationship between data found no relationship. If you are creating data with related data id, those id's do not exist."
	default:
		statuscode = codes.Internal
		message = "internal error"
	}

	return status.Error(statuscode, message)
}
