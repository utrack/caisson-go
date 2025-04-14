package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/errorbag"
)

func SetLevel(level slog.Level) {
	slog.SetLogLoggerLevel(level)
}

func Debug(ctx context.Context, msg string, kvs ...any) {
	fromCtx(ctx).DebugContext(ctx, msg, kvs...)
}

func Info(ctx context.Context, msg string, kvs ...any) {
	fromCtx(ctx).InfoContext(ctx, msg, kvs...)
}

func Warn(ctx context.Context, msg string, kvs ...any) {
	fromCtx(ctx).WarnContext(ctx, msg, kvs...)
}

func Warne(ctx context.Context, msg string, err error, kvs ...any) {
	if err != nil {
		Warn(ctx, msg, errKvs(err, kvs)...)
	}
}

func errKvs(err error, kvs []any) []any {
	kvs = append(kvs, "error.message", err.Error())
	kvs = append(kvs, "error.stack", fmt.Sprintf("%+v", err))
	if code := errors.Code(err); code != nil {
		kvs = append(kvs, "error.code", code.Type())
		kvs = append(kvs, "error.user_message", code.Message())
	}

	data := errorbag.ListPairs(err)
	if len(data) > 0 {
		kvs = append(kvs, "error.data", data)
	}
	return kvs
}

func Error(ctx context.Context, msg string, err error, kvs ...any) {
	if err != nil {
		Errorn(ctx, msg, errKvs(err, kvs)...)
	}
}

// Errorn emits a log with level Error but it does not add any error context.
// In 99% of the cases you want to use Error or Errorne instead.
func Errorn(ctx context.Context, msg string, kvs ...any) {
	fromCtx(ctx).ErrorContext(ctx, msg, kvs...)
}

// Errorne is a shortcut for Error(ctx, err.Error(), err, errKvs(err, kvs)...)
// It's useful when you want to log an error without any additional context.
func Errorne(ctx context.Context, err error, kvs ...any) {
	if err != nil {
		Error(ctx, err.Error(), err, kvs...)
	}
}

func Fatal(ctx context.Context, msg string, kvs ...any) {
	fromCtx(ctx).ErrorContext(ctx, msg, kvs...)
	os.Exit(1)
}

func Fatale(ctx context.Context, msg string, err error, kvs ...any) {
	if err != nil {
		Fatal(ctx, msg, errKvs(err, kvs)...)
	}
}
