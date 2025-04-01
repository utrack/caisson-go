package errors

import (
	"fmt"

	"github.com/utrack/caisson-go/pkg/errorbag"
)

// WithKeyedData adds a typed detail to the error.
// Duplicate keys overlay each other, but it's still possible to retrieve
// all values associated with the key via errorbag.GetAll().
//
// You can retrieve a keyed value later via KeyedData().
//
// Nil cause returns a nil error.
//
// When used in conjunction with errors.As(),
// errors.As(err, &T{}) will return the first value of the type.
func WithKeyedData[K comparable, T any](err error, key K, value T) error {
	return errorbag.With(err, key, value)
}

// KeyedData retrieves the value associated with the given key.
//
// Returns (zero value, false) if the key is not found.
func KeyedData[K comparable, T any](err error, key K) (T, bool) {
	return errorbag.Get[K, T](err, key)
}

func Wrapd(cause error, msg string, kvs ...any) error {
	cause = Wrap(cause, msg)
	for i := 0; i < len(kvs); i += 2 {
		if i+1 >= len(kvs) {
			break
		}
		k := fmt.Sprintf("%v", kvs[i])
		v := kvs[i+1]
		cause = WithKeyedData(cause, k, v)
	}

	return cause
}
