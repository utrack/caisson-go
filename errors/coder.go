package errors

import (
	"fmt"
)

// Coder is a static struct which enriches passing errors with HTTP codes and user messages.
type Coder interface {
	WithMessage(userMessage string) Coder
	WithMessagef(format string, args ...any) Coder
	WithHTTPCode(httpCode int) Coder
	Wrap(cause error) error
}

// Coded is an interface for errors enriched with HTTP codes and user messages.
type Coded interface {
	HTTPCode() int
	Message() string
	Cause() error
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

type coder struct {
	httpCode    int
	userMessage string
	cause       error
}

var _ Coder = coder{}

func (c coder) WithMessage(userMessage string) Coder {
	return coder{
		httpCode:    c.httpCode,
		userMessage: userMessage,
		cause:       c.cause,
	}
}

func (c coder) WithMessagef(format string, args ...any) Coder {
	return coder{
		httpCode:    c.httpCode,
		userMessage: fmt.Sprintf(format, args...),
		cause:       c.cause,
	}
}

func (c coder) WithHTTPCode(httpCode int) Coder {
	return coder{
		httpCode:    httpCode,
		userMessage: c.userMessage,
		cause:       c.cause,
	}
}

func (c coder) Wrap(cause error) error {
	return coded{
		httpCode:    c.httpCode,
		userMessage: c.userMessage,
		cause:       cause,
	}
}

type coded struct {
	httpCode    int
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

func (c coded) Cause() error {
	return c.cause
}

func (c coded) Error() string {
	return c.cause.Error()
}
