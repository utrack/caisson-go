package httpbinding

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/ggicci/httpin"
	"github.com/utrack/caisson-go/errors"
)

// requestDecoder decodes a request body into a type T.
type requestDecoder interface {
	Decode(req *http.Request) (any, error)
}

// jsonBodyDecoder decodes a request body using JSON body only.
type jsonBodyDecoder struct {
	t reflect.Type
}

func newJsonBodyDecoder(t reflect.Type) requestDecoder {
	return jsonBodyDecoder{t: t}
}

func (j jsonBodyDecoder) Decode(req *http.Request) (any, error) {
	out := reflect.New(j.t).Interface()
	if err := json.NewDecoder(req.Body).Decode(out); err != nil {
		return nil, errors.Wrap(err, "failed to decode JSON request body")
	}
	return out, nil
}

func newHttpinDecoder(t reflect.Type) (requestDecoder, error) {
	inValue := reflect.New(t).Interface()

	unmEngine, err := httpin.New(inValue)
	if err != nil {
		return nil, err
	}

	return unmEngine, nil
}

func hasTag(t reflect.Type, tag string) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if _, ok := field.Tag.Lookup(tag); ok {
			return true
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if hasTag(field.Type, tag) {
				return true
			}
		}
	}
	return false
}
