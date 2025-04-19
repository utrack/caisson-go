package handler

import (
	"net/http"

	"github.com/utrack/caisson-go/pkg/http/hhandler"
)

type OptionHTTP = hhandler.OptionHTTP

func WithGlobalMiddleware(middleware func(http.Handler) http.Handler) OptionHTTP {
	return func(o *hhandler.Options) {
		o.Middlewares = append(o.Middlewares, middleware)
	}
}