package errors

import (
	"fmt"
)

// Coder is a static struct which enriches passing errors with HTTP codes and user messages.
type Coder interface {
	WithType(typ string) Coder
	WithMessage(userMessage string) Coder
	WithMessagef(format string, args ...any) Coder
	WithHTTPCode(httpCode int) Coder
	Wrap(cause error) error
}

// Coded is an interface for errors enriched with HTTP codes and user messages.
type Coded interface {
	HTTPCode() int
	Message() string
	Type() string
	Unwrap() error
}

type keyCoderType string

const (
	keyCoded keyCoderType = "errors.Coded"
)

func Code(err error) Coded {
	c, ok := Data[keyCoderType, Coded](err, keyCoded)
	if ok {
		return c
	}
	return nil
}

func NewCoder(typ string) Coder {
	return coder{typ: typ}
}

type coder struct {
	httpCode    int
	typ         string
	userMessage string
}

var _ Coder = coder{}

func (c coder) WithType(typ string) Coder {
	return coder{
		httpCode:    c.httpCode,
		typ:         typ,
		userMessage: c.userMessage,
	}
}

func (c coder) WithMessage(userMessage string) Coder {
	return coder{
		httpCode:    c.httpCode,
		typ:         c.typ,
		userMessage: userMessage,
	}
}

func (c coder) WithMessagef(format string, args ...any) Coder {
	return coder{
		httpCode:    c.httpCode,
		typ:         c.typ,
		userMessage: fmt.Sprintf(format, args...),
	}
}

func (c coder) WithHTTPCode(httpCode int) Coder {
	return coder{
		httpCode:    httpCode,
		typ:         c.typ,
		userMessage: c.userMessage,
	}
}

func (c coder) Wrap(cause error) error {
	return coded{
		httpCode:    c.httpCode,
		typ:         c.typ,
		userMessage: c.userMessage,
		cause:       cause,
	}
}

type coded struct {
	httpCode    int
	typ         string
	userMessage string
	cause       error
}

var _ Coded = coded{}

func (c coded) HTTPCode() int {
	if c.httpCode == 0 {
		return 500
	}
	return c.httpCode
}

func (c coded) Message() string {
	inner := Code(c.cause)

	if inner != nil && inner.Message() != "" {
		return inner.Message() + ": " + c.userMessage
	}

	return c.userMessage
}

func (c coded) Type() string {
	return c.typ
}

func (c coded) Unwrap() error {
	return c.cause
}

func (c coded) Error() string {
	return c.cause.Error()
}
