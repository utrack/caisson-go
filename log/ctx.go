package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type ctxLogKeyType string

const ctxLogKey ctxLogKeyType = "caisson.logger"

func With(ctx context.Context, kvs ...any) context.Context {
	current := ctxLogger(ctx)
	newLogger := current.With(kvs...)
	return context.WithValue(ctx, ctxLogKey, newLogger)
}

func ctxLogger(ctx context.Context) *slog.Logger {
	ret, ok := ctx.Value(ctxLogKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return ret
}

func fromCtx(ctx context.Context) *slog.Logger {
	l := ctxLogger(ctx)

	// extraction taken from slogotel "github.com/veqryn/slog-context/otel", but it does more than what we need
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		spanCtx := span.SpanContext()
		if spanCtx.HasTraceID() {
			l = l.With(slog.String("trace_id", spanCtx.TraceID().String()))
		}
		if spanCtx.HasSpanID() {
			l = l.With(slog.String("span_id", spanCtx.SpanID().String()))
		}
	}
	return l
}
