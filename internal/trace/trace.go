package trace

import (
	"context"
	"fmt"
	"io"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func SetOTLPGRPCTracerProvider(ctx context.Context) (func(context.Context) error, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("creating new otlp trace grpc exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(NewResource(ctx)),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

func SetOTLPHTTPTracerProvider(ctx context.Context) (func(context.Context) error, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating new otlp trace grpc exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(NewResource(ctx)),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

// NewStdOutExporter returns a console exporter.
func NewStdOutExporter(ctx context.Context, w io.Writer) (func(context.Context) error, error) {
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new std out exporter")
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(NewResource(ctx)),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

// NewResource returns a resource describing this application.
func NewResource(ctx context.Context) *resource.Resource {
	r, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(meta.ServiceName),
			semconv.ServiceVersionKey.String("v0.0.0"),
			attribute.String("environment", "local"),
		),
	)
	return r
}
