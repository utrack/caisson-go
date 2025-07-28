package logctx

import (
	"context"
	"log/slog"
)

type ctxLogKeyType string

const ctxLogKey ctxLogKeyType = "caisson.l3.logctx"

// With returns a new context with the logger enriched with the given key-value pairs.
// Any calls to the logger from the new context will include the given key-value pairs.
func With(ctx context.Context, kvs ...any) context.Context {
	current := From(ctx)
	newLogger := current.With(kvs...)
	return context.WithValue(ctx, ctxLogKey, newLogger)
}

// From returns the logger from the context.
// If the context does not have a logger, returns the default [slog.Default()] logger.
func From(ctx context.Context) *slog.Logger {
	ret, ok := ctx.Value(ctxLogKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return ret
}
