package pgx

import (
	"context"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct{}

func (t *Tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "pgx.TraceQueryStart", trace.WithAttributes(
		attribute.String("sql", data.SQL),
	))
	defer span.End()
	return ctx
}

func (t *Tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	_, span := otel.Tracer(meta.ServiceName).Start(ctx, "pgx.TraceQueryEnd", trace.WithAttributes(
		attribute.String("sql", data.CommandTag.String()),
	))
	defer span.End()
}
