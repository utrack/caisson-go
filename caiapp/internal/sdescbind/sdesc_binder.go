package sdescbind

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/http/errmarshalhttp"
	"github.com/utrack/caisson-go/pkg/http/hhandler"
	"github.com/utrack/caisson-go/pkg/http/httpbinding"
	"github.com/utrack/caisson-go/pkg/http/negmarshal"
	"github.com/utrack/pontoon/sdesc"
)

func Bind(s sdesc.Service, h hhandler.Handler) error {
	sconfig := sdesc.HandlerConfig{}
	for _, opt := range s.ServiceOptions() {
		opt(&sconfig)
	}
	b := &binder{
		errorRender: func(ctx context.Context, r *http.Request, w http.ResponseWriter, err error) {
			errmarshalhttp.Marshal(err, w, r)
		},
		neg:     negmarshal.Default(),
		h:       h,
		sconfig: sconfig,
	}
	s.RegisterHTTP(b)

	return b.bindError
}

type binder struct {
	errorRender httpbinding.ErrorRenderer
	neg         negmarshal.NegotiatedMarshalFunc
	h           hhandler.Handler
	sconfig     sdesc.HandlerConfig

	bindError error
}

var _ sdesc.HTTPRouter = (*binder)(nil)

func (b *binder) MethodFunc(method, pattern string, hdl sdesc.RPCHandler) {
	handler, err := httpbinding.BindHTTPHandler(hdl, b.errorRender, b.neg)
	if err != nil {
		b.bindError = errors.Wrapd(err, "when binding HTTP handler", "method", method, "pattern", pattern)
		return
	}

	mws := chi.Middlewares(b.sconfig.Middlewares())

	b.h.MethodFunc(method, pattern, mws.HandlerFunc(handler.ServeHTTP).ServeHTTP)
}
