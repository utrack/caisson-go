package errors

import "github.com/utrack/caisson-go/pkg/errorbag"

func WithData[K comparable, T any](err error, key K, value T) error {
	return errorbag.With(err, key, value)
}

func Data[K comparable, T any](err error, key K) (T, bool) {
	return errorbag.Get[K, T](err, key)
}
