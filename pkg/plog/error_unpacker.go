package plog

import (
	"log/slog"

	"github.com/utrack/caisson-go/levels/level3/errorbag"
)

// errorUnpacker unpacks an error into a slog.Value.
//
// It extracts any key-value pairs associated with the error,
// embedded via either errors or pkg/errorbag.
type errorUnpacker struct {
	err error
}

func errLogged(err error) errorUnpacker {
	return errorUnpacker{err: err}
}

func (u errorUnpacker) LogValue() slog.Value {
	kvs := errorbag.ListPairs(u.err)
	return slog.AnyValue(kvs)
}
