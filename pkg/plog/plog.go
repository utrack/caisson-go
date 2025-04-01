package plog

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

type ctxLogKeyType string

const ctxLogKey ctxLogKeyType = "plogKey"

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

func Info(ctx context.Context, msg string, kvs ...any) {
	current := fromCtx(ctx)
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, pcs[0])
	r.Add(kvs...)
	_ = current.Handler().Handle(ctx, r)
}

func Error(ctx context.Context, err error, kvs ...any) {
	current := fromCtx(ctx)
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])

	r := slog.NewRecord(time.Now(), slog.LevelError, err.Error(), pcs[0])
	kvs = append(kvs, "error", errLogged(err))
	r.Add(kvs...)
	_ = current.Handler().Handle(ctx, r)
}
