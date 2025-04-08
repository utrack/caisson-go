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

func Code(err error) Coded {
	if err == nil {
		return nil
	}
	var ret Detailed[Coded]
	ok := As(err, &ret)
	if ok {
		return ret.Details()
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
	return DetailWith[Coded](cause, coded{
		HttpCode:    c.httpCode,
		Typ:         c.typ,
		UserMessage: c.userMessage,
		cause:       cause,
	})
}

type coded struct {
	HttpCode    int `json:"http_code"`
	Typ         string `json:"type"`
	UserMessage string `json:"user_message"`
	cause       error
}

var _ Coded = coded{}

func (c coded) HTTPCode() int {
	if c.HttpCode == 0 {
		return 500
	}
	return c.HttpCode
}

func (c coded) Message() string {
	inner := Code(c.cause)

	if inner != nil && inner.Message() != "" {
		return c.UserMessage + ": " + inner.Message()
	}

	return c.UserMessage
}

func (c coded) Type() string {
	return c.Typ
}

func (c coded) Unwrap() error {
	return c.cause
}

func (c coded) Error() string {
	return c.cause.Error()
}
