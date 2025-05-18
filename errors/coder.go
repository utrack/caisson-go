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

// Coded is an error instance enriched with HTTP codes and user messages.
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
	return coder{data: coded{
		Typ:         typ,
		HttpCode:    500,
		UserMessage: "",
	}}
}

type coder struct {
	data coded
}

var _ Coder = coder{}

func (c coder) WithType(typ string) Coder {
	d := c.data
	d.Typ = typ
	return coder{
		data: d,
	}
}

func (c coder) WithMessage(userMessage string) Coder {
	d := c.data
	d.UserMessage = userMessage
	return coder{
		data: d,
	}
}

func (c coder) WithMessagef(format string, args ...any) Coder {
	d := c.data
	d.UserMessage = fmt.Sprintf(format, args...)
	return coder{
		data: d,
	}
}

func (c coder) WithHTTPCode(httpCode int) Coder {
	d := c.data
	d.HttpCode = httpCode
	return coder{
		data: d,
	}
}

func (c coder) Wrap(cause error) error {
	d := c.data
	d.cause = cause
	return DetailWith[Coded](cause, d)
}

type coded struct {
	HttpCode    int    `json:"http_code"`
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
