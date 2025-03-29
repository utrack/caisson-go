package errors

import (
	std "errors"

	pkg "github.com/pkg/errors"
)

func New(text string) error {
	// std.New() does not add a stack trace
	return pkg.New(text)
}

func Errorf(format string, args ...any) error {
	// std.Errorf() does not add a stack trace
	return pkg.Errorf(format, args...)
}

func Is(err error, target error) bool {
	return std.Is(err, target)
}

func As(err error, target any) bool {
	return std.As(err, target)
}

func Unwrap(err error) error {
	return std.Unwrap(err)
}

func Join(errs ...error) error {
	return std.Join(errs...)
}
