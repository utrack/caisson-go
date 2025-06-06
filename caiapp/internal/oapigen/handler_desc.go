package oapigen

import (
	"reflect"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/utrack/pontoon/v2/httpinoapi"
)

type HandlerDesc struct {
	Method string
	Path   string

	Func   any
	Input  reflect.Type
	Output reflect.Type
}

func GenerateOAPI(handlers []HandlerDesc) (*v3.Document, error) {
	gen := httpinoapi.NewGenerator()

	for _, d := range handlers {
		opts := []httpinoapi.Option{}
		if d.Input != nil {
			st := reflect.New(d.Input).Elem().Interface()
			opts = append(opts, httpinoapi.WithInputStruct(st))
		}
		if d.Output != nil {
			st := reflect.New(d.Output).Elem().Interface()
			opts = append(opts, httpinoapi.WithOutputStruct(st))
		}
		gen.Operation(d.Method, d.Path, d.Func, opts...)
	}

	return gen.Build()
}
