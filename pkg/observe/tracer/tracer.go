package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan creates a new child span from the parent (if it exists).
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer("").Start(ctx, name)
}
