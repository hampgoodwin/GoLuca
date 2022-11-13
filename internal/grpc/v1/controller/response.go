package controller

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/hampgoodwin/errors"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) respondError(ctx context.Context, log *zap.Logger, err error) error {
	_, span := otel.Tracer(meta.ServiceName).Start(ctx, "http.v0.controller.respondError")
	span.RecordError(err)
	log.Error("responding", zap.Error(err))

	var statuscode codes.Code
	var message string
	var msg errors.Message
	if errors.As(err, &msg) {
		message = msg.Value
	}
	span.SetStatus(otelcodes.Error, message)

	switch {
	case errors.Is(err, errors.NotKnown):
		statuscode = codes.Unknown
	case errors.Is(err, errors.NotValidRequest):
		statuscode = codes.InvalidArgument
		message = "bad request data, check request meta data."
	case errors.Is(err, errors.NotValidRequestData):
		statuscode = codes.InvalidArgument
		if message == "" {
			message = "bad request data"
		}
		if err := c.respondOnValidationErrors(err, message); err != nil {
			return err
		}
		c.log.Error(
			"incorrect error flag used for case",
			zap.Error(err), zap.String("error_flag", fmt.Sprint(errors.NotValidRequestData)),
		)
	case errors.Is(err, errors.NotFound):
		statuscode = codes.NotFound
	case errors.Is(err, errors.NotValidInternalData):
		statuscode = codes.Internal
		message = "internal data is invalid and failed validation."
	case errors.Is(err, errors.NotDeserializable):
		statuscode = codes.Internal
		message = "provided data passed failed deserialization. If creating a resource, check the request body types."
	case errors.Is(err, errors.NotSerializable):
		statuscode = codes.Internal
		message = "either provided or internal data passed validation, but failed serialization."
	case errors.Is(err, errors.NoRelationshipFound):
		statuscode = codes.InvalidArgument
		message = "process which assumed existence of a relationship between data found no relationship. If you are creating data with related data id, those id's do not exist."
	default:
		statuscode = codes.Internal
		message = "internal error"
	}
	if errors.As(err, &msg) {
		if msg.Value != "" {
			message = msg.Value
		}
	}

	return status.Error(statuscode, message)
}

func (c *Controller) respondOnValidationErrors(err error, message string) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		st := status.New(codes.InvalidArgument, validationErrors.Error())
		br := &errdetails.BadRequest{}
		for _, fieldError := range validationErrors {
			fv := &errdetails.BadRequest_FieldViolation{
				Field:       fieldError.StructNamespace(),
				Description: fieldError.Tag(),
			}
			br.FieldViolations = append(br.FieldViolations, fv)
		}
		st, err := st.WithDetails(br)
		if err != nil {
			return status.Error(codes.Internal, "error construction validation error response")
		}
		return st.Err()
	}
	return nil
}
