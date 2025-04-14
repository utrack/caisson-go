package slogtrace

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

var _ slog.Handler = &contextAdapter{}

type contextAdapter struct {
	inner slog.Handler
}

// NewContextHandler returns a new [slog.Handler] that adds trace and span IDs to log records.
func NewContextHandler(inner slog.Handler) slog.Handler {
	return &contextAdapter{inner: inner}
}

// Handle implements [slog.Handler].
func (c *contextAdapter) Handle(ctx context.Context, r slog.Record) error {
	levelText := r.Level.String()
	if r.Level > 20 { // according to spec, https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-severitytext
		levelText = "FATAL"
	}
	r.Add(slog.String("severity_text", levelText))
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		spanCtx := span.SpanContext()
		if spanCtx.HasTraceID() {
			r.Add(slog.String("trace_id", spanCtx.TraceID().String()))
		}
		if spanCtx.HasSpanID() {
			r.Add(slog.String("span_id", spanCtx.SpanID().String()))
		}
	}

	return c.inner.Handle(ctx, r)
}

// WithAttrs implements [slog.Handler].
func (c *contextAdapter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextAdapter{c.inner.WithAttrs(attrs)}
}

// WithGroup implements [slog.Handler].
func (c *contextAdapter) WithGroup(name string) slog.Handler {
	return &contextAdapter{c.inner.WithGroup(name)}
}

// Enabled implements [slog.Handler].
func (c *contextAdapter) Enabled(ctx context.Context, level slog.Level) bool {
	return c.inner.Enabled(ctx, level)
}
