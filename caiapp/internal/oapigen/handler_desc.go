package oapigen

import (
	"path"
	"reflect"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/utrack/caisson-go/caiapp/internal/hchi"
	"github.com/utrack/pontoon/v2/httpinoapi"
)

type HandlerDesc struct {
	Method string
	Path   string

	Func   any
	Input  reflect.Type
	Output reflect.Type
}

func GenerateOAPI(handlers []HandlerDesc, ropts hchi.OptionExtensions) (*v3.Document, error) {
	gen := httpinoapi.NewGenerator()

	for _, d := range handlers {
		opts := []httpinoapi.Option{}
		if d.Input != nil {
			opts = append(opts, httpinoapi.WithInputType(d.Input))
		}
		if d.Output != nil {
			opts = append(opts, httpinoapi.WithOutputType(d.Output))
		}
		gen.Operation(d.Method, path.Join(ropts.Prefix, d.Path), d.Func, opts...)
	}

	return gen.Build()
}
