package errors

import (
	"reflect"

	"github.com/utrack/caisson-go/levels/level3/errorbag"
)

// DetailWith adds a typed detail to the error.
// To retrieve it, cast the error to Detailed[T] and call Details() like so:
//
// var d Detailed[MyDetailType]
// ok := errors.As(err, &d)
//
//	if ok {
//	    details := d.Value()
//	}
func DetailWith[T any](err error, value T) error {
	if err == nil {
		return nil
	}
	typeName := reflect.TypeFor[T]()
	return WithKeyedData(err, typeName, value)
}

// Detailed is an interface for errors enriched with typed details.
type Detailed[T any] interface {
	error
	Value() T
}

var _ Detailed[any] = errorbag.Bag[string, any](nil)
