package errors

import pkg "github.com/pkg/errors"

func Cause(err error) error {
	return pkg.Cause(err)
}

func WithMessage(err error, message string) error {
	return pkg.WithMessage(err, message)
}

func WithMessagef(err error, format string, args ...any) error {
	return pkg.WithMessagef(err, format, args...)
}

func WithStack(err error) error {
	return pkg.WithStack(err)
}

func Wrap(err error, message string) error {
	return pkg.Wrap(err, message)
}

func Wrapf(err error, format string, args ...any) error {
	return pkg.Wrapf(err, format, args...)
}
