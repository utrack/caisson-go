package errors

import (
	"errors"
	"fmt"
)

// Coder is a static struct which enriches passing errors with HTTP codes and user messages.
type Coder interface {
	WithType(typ string) Coder
	WithMessage(userMessage string) Coder
	WithMessagef(format string, args ...any) Coder
	WithHTTPCode(httpCode int) Coder
	Wrap(cause error) error
	Error() string
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
		return ret.Value()
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

func (c coder) Error() string {
	if c.data.Typ == "" {
		return c.data.UserMessage
	}
	return fmt.Sprintf("%s: %s", c.data.Typ, c.data.UserMessage)
}

func (c coder) Is(target error) bool {
	var t coder
	// TODO test-cover
	return errors.As(target, &t) && t.data.Typ == c.data.Typ && t.data.HttpCode == c.data.HttpCode && t.data.UserMessage == c.data.UserMessage
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

func (c coded) Is(target error) bool {
	var t Coded
	return errors.As(target, &t) && t.Type() == c.Typ && t.HTTPCode() == c.HttpCode && t.Message() == c.UserMessage
}
