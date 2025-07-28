// Package errorbag provides a way to associate
// errors with arbitrary keys and values.
//
// Bags are used to add context to errors - as an example, you can use them to pass
// HTTP/gRPC codes, error reasons, aux error data for the frontend, etc.
// They work as a stack of key-value pairs; like an inverted `context` value, where the stack
// is populated from the innermost code part (the error cause) instead of the outermost one (API request).
//
// Bags implement the error interface and can be
// used with the errors, github.com/pkg/errors and other packages.
package errorbag

import (
	"errors"
	"fmt"
)

// Bag is a key-value pair with an error
// associated with it.
type Bag[K comparable, T any] interface {
	error
	Key() K
	Value() T
	Unwrap() error
	bagAny
}

type bagAny interface {
	keyAny() any
	valueAny() any
}

// With creates a new bag with the given cause, key, and value.
// You can wrap the same error chain many times to add many K/V pairs.
//
// Nil cause returns a nil Bag.
//
// The 'context' key rules apply; it's better to create your own key type and use it instead of a string to protect against collisions.
func With[K comparable, T any](cause error, key K, value T) Bag[K, T] {
	if cause == nil {
		return nil
	}
	return container[K, T]{
		cause: cause,
		key:   key,
		value: value,
	}
}

// Get returns the value associated with the given key; recursing into the error chain to find it.
//
// It follows the errors.Unwrap() semantics; key-value pairs behind the errors.Join are unsupported.
func Get[K comparable, T any](err error, key K) (T, bool) {
	for err != nil {
		if bag, ok := err.(Bag[K, T]); ok {
			if bag.Key() == key {
				return bag.Value(), true
			}
		}
		err = errors.Unwrap(err)
	}
	return *new(T), false
}

// GetAll returns all the values associated with the given key; recursing into the error chain to find them.
//
// It follows the errors.Unwrap() semantics; key-value pairs behind the errors.Join are unsupported.
func GetAll[K comparable, T any](err error, key K) ([]T, bool) {
	var ret []T
	for err != nil {
		if bag, ok := err.(Bag[K, T]); ok {
			if bag.Key() == key {
				ret = append(ret, bag.Value())
			}
		}
		err = errors.Unwrap(err)
	}
	return ret, len(ret) > 0
}

// ListPairs returns all the key-value pairs associated with the given error; recursing into the error chain to find them.
//
// Listing the pairs is a potentially expensive operation, as it requires
// traversing the entire error chain.
func ListPairs(err error) map[string]any {
	ret := make(map[string]any, 3)
	for err != nil {
		if bag, ok := err.(bagAny); ok {
			keyStr := fmt.Sprintf("%s", bag.keyAny())
			if v, ok := ret[keyStr]; ok {
				if asArray, ok := v.([]any); ok {
					asArray = append(asArray, bag.valueAny())
					ret[keyStr] = asArray
				} else {
					ret[keyStr] = []any{v, bag.valueAny()}
				}
			} else {
				ret[keyStr] = bag.valueAny()
			}
		}
		err = errors.Unwrap(err)
	}
	return ret
}

type container[K comparable, T any] struct {
	cause error

	key   K
	value T
}

func (c container[K, T]) Error() string {
	return c.cause.Error()
}

func (c container[K, T]) Key() K {
	return c.key
}

func (c container[K, T]) Value() T {
	return c.value
}

func (c container[K, T]) Unwrap() error {
	return c.cause
}

func (c container[K, T]) keyAny() any {
	return c.key
}

func (c container[K, T]) valueAny() any {
	return c.value
}
