package errors

// CoderDetailer is a Coder that enriches errors with typed details.
//
// See Coder for the general description.
type CoderDetailer[T any] interface {
	WithType(typ string) CoderDetailer[T]
	WithMessage(userMessage string) CoderDetailer[T]
	WithMessagef(format string, args ...any) CoderDetailer[T]
	WithHTTPCode(httpCode int) CoderDetailer[T]
	Wrap(cause error, details T) error
	
	// Details extracts the embedded details from the error decorated via Wrap().
	Details(err error) *T
}

func NewCoderDetailer[T any](typ string) CoderDetailer[T] {
	return coderDetailer[T]{coder: NewCoder(typ)}
}

type coderDetailer[T any] struct {
	coder Coder
}

var _ CoderDetailer[string] = coderDetailer[string]{}

func (c coderDetailer[T]) WithType(typ string) CoderDetailer[T] {
	return coderDetailer[T]{
		coder: c.coder.WithType(typ),
	}
}

func (c coderDetailer[T]) WithMessage(userMessage string) CoderDetailer[T] {
	return coderDetailer[T]{
		coder: c.coder.WithMessage(userMessage),
	}
}

func (c coderDetailer[T]) WithMessagef(format string, args ...any) CoderDetailer[T] {
	return coderDetailer[T]{
		coder: c.coder.WithMessagef(format, args...),
	}
}

func (c coderDetailer[T]) WithHTTPCode(httpCode int) CoderDetailer[T] {
	return coderDetailer[T]{
		c.coder.WithHTTPCode(httpCode),
	}
}

func (c coderDetailer[T]) Wrap(cause error, details T) error {
	// TODO use error type instead of T's reflect type
	return DetailWith(c.coder.Wrap(cause), details)
}

func (c coderDetailer[T]) Details(err error) *T {
	var d Detailed[T]
	ok := As(err, &d)
	if ok {
		ret := d.Details()
		return &ret
	}
	return nil
}