package httpbinding

import (
	"context"
	"net/http"
	"reflect"

	"github.com/ggicci/httpin/core"
	"github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/http/negmarshal"
	"github.com/utrack/pontoon/sdesc"
)

var (
	errorInterface  = reflect.TypeOf((*error)(nil)).Elem()
	writerInterface = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	typeHttpReq     = reflect.TypeOf(&http.Request{})
	typeCtx         = reflect.TypeOf(context.Context(nil))
)

func init() {
	integration.UseGochiURLParam("path", chi.URLParam)
}

var ErrMalformedRequest = errors.NewCoder("BAD_REQUEST").WithHTTPCode(400).WithMessage("malformed request body")

type ErrorRenderer func(context.Context, *http.Request, http.ResponseWriter, error)

// wrapDescRPCHandler converts sdesc.RPCHandler to stdlib http.HandlerFunc.
// It can wrap handlers that accept any/all of *http.Request, http.ResponseWriter
// and any custom struct (which is then unmarshaled via ggicci/httpin).
//
// Handlers' output types are either ({return type},error), (error) or nothing.
// The return type is marshaled to negotiated content type, or JSON by default.
//
// Writing to http.ResponseWriter is not allowed if handler has a return type.
func BindHTTPHandler(h sdesc.RPCHandler, errRender ErrorRenderer, marshaler negmarshal.NegotiatedMarshalFunc) (http.Handler, error) {
	ret, _, err := BindHTTPHandlerMeta(h, errRender, marshaler)
	return ret, err
}

// BindHTTPHandlerMeta is like BindHTTPHandler, but returns the handler's meta information.
//
// It can be used to inspect the handler's input and output types for documentation autogen.
//
// TODO this can be split into two functions, one for meta extraction and one for binding based on the meta.
func BindHTTPHandlerMeta(h sdesc.RPCHandler, errRender ErrorRenderer, marshaler negmarshal.NegotiatedMarshalFunc) (http.Handler, Meta, error) {
	if errRender == nil {
		return nil, Meta{}, errors.New("nil ErrorRenderer")
	}
	handleFuncRef := reflect.ValueOf(h)
	if handleFuncRef.Kind() != reflect.Func {
		return nil, Meta{}, errors.New("handler is not a function")
	}
	funcType := handleFuncRef.Type()

	if funcType.NumIn() > 3 {
		return nil, Meta{}, errors.New("handler should accept 1 to 3 parameters")
	}

	if funcType.NumOut() == 2 && !funcType.Out(1).Implements(errorInterface) {
		return nil, Meta{}, errors.New("2nd return type should be an error")
	}

	hasOutputStruct := funcType.NumOut() > 1

	// if a function accepts a ResponseWriter - assume it's going to write to it directly;
	// we don't control the output anymore.
	var controlsResponseWriter bool

	// functions that convert i-th input type of a handler function
	// to reflect.Value() for calling
	type inFun func(w http.ResponseWriter, r *http.Request) (reflect.Value, error)

	inFuncs := []inFun{}
	var inType reflect.Type
	for i := 0; i < funcType.NumIn(); i++ {
		switch {
		case funcType.In(i) == typeHttpReq:
			inFuncs = append(inFuncs, func(_ http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(r), nil
			})
		case funcType.In(i).Implements(writerInterface):
			controlsResponseWriter = true
			inFuncs = append(inFuncs, func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(w), nil
			})
		default:
			inType = funcType.In(i)

			unmEngine, err := newHttpinDecoder(inType)
			if err != nil {
				return nil, Meta{}, errors.Wrapf(err, "failed to create HTTPin decoder for type %v %v", inType.PkgPath(), inType.Name())
			}

			inFuncs = append(inFuncs, func(_ http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				in, err := unmEngine.Decode(r)
				if err != nil {
					var invalidFieldError *core.InvalidFieldError
					if errors.As(err, &invalidFieldError) {
						return reflect.Value{}, ErrMalformedRequest.Wrap(err)
					}
					return reflect.Value{}, errors.Wrap(err, "failed to decode HTTPin request")
				}
				return reflect.ValueOf(in).Elem(), nil
			})
		}
	}
	if funcType.NumOut() < 1 && !controlsResponseWriter {
		return nil, Meta{}, errors.New("handler should return error if it doesn't accept http.ResponseWriter")
	}
	if funcType.NumOut() > 2 {
		return nil, Meta{}, errors.New("handler should return maximum of 2 parameters")
	}
	if funcType.NumOut() > 0 && controlsResponseWriter {
		return nil, Meta{}, errors.New("handler should not return anything if it controls the response writer directly")
	}

	retMeta := Meta{
		NamedFunc:         h,
		WriterIntercepted: controlsResponseWriter,
	}

	if funcType.NumOut() > 1 {
		retMeta.OutputType = funcType.Out(0)
	}

	if inType != nil {
		retMeta.InputType = inType
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inArgs := []reflect.Value{}

		if !controlsResponseWriter {
			ow := &wrappedWriter{blockWrites: hasOutputStruct, w: w}
			w = ow
		}
		for _, f := range inFuncs {
			v, err := f(w, r)
			if err != nil {
				errRender(r.Context(), r, w, err)
				return
			}
			inArgs = append(inArgs, v)
		}

		out := handleFuncRef.Call(inArgs)
		if len(out) == 0 {
			return
		}

		errPos := 1
		if !hasOutputStruct {
			errPos = 0
		}

		if !out[errPos].IsNil() {
			errRender(r.Context(), r, w, out[errPos].Interface().(error))
			return
		}

		var err error
		switch {
		case hasOutputStruct:
			err = marshaler(r, w, out[0].Interface())
		case !controlsResponseWriter:
			// do not output anything if the handler controls the writer directly
		default:
			err = marshaler(r, w, struct{}{})
		}
		if err != nil {
			err = errors.Wrap(err, "the call succeeded, but failed to marshal the response")
			errRender(r.Context(), r, w, err)
		}
	}), retMeta, nil
}
