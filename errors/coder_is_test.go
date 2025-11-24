package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoderIs(t *testing.T) {
	so := require.New(t)

	code := NewCoder("test").WithHTTPCode(400).WithMessage("base Coder")
	code2 := NewCoder("test2").WithHTTPCode(400).WithMessage("base Coder")
	so.True(Is(code, code))
	so.False(Is(code, code2))

	someErr := errors.New("some basic error")
	so.False(Is(code, someErr))

	wr := code.Wrap(someErr)
	so.True(Is(wr, someErr))
	so.True(Is(wr, code))
	so.False(Is(wr, code2))
}
