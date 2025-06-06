package httpbinding

import "reflect"

// Meta is the metadata of a bound HTTP handler.
// It is used for later documentation generation.
type Meta struct {
	InputType  reflect.Type
	OutputType reflect.Type

	NamedFunc any

	WriterIntercepted bool
}
