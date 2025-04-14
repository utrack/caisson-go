package log

import (
	"context"
	"log/slog"
)

type ctxLogKeyType string

const ctxLogKey ctxLogKeyType = "caisson.logger"

func With(ctx context.Context, kvs ...any) context.Context {
	current := fromCtx(ctx)
	newLogger := current.With(kvs...)
	return context.WithValue(ctx, ctxLogKey, newLogger)
}

func fromCtx(ctx context.Context) *slog.Logger {
	ret, ok := ctx.Value(ctxLogKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return ret
}
