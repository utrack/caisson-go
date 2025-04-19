package hchi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/utrack/caisson-go/pkg/http/hhandler"
)

var _ hhandler.Handler = &ChiHandler{}

type ChiHandler struct {
	// we use delayed registration, so that we can apply any middlewares before the actual route registration.
	routes  []route
	options hhandler.Options
}

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

func New() *ChiHandler {
	srv := &http.Server{
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Minute * 2,
	}
	return &ChiHandler{
		options: hhandler.Options{
			Server: srv,
		},
	}
}

func (c *ChiHandler) Apply(oo ...hhandler.OptionHTTP) {
	for _, o := range oo {
		o(&c.options)
	}
}

func (c *ChiHandler) MethodFunc(method string, pattern string, hdl http.HandlerFunc) {
	c.routes = append(c.routes, route{method: method, pattern: pattern, handler: hdl})
}

func (c *ChiHandler) HandleFunc(pattern string, hdl http.HandlerFunc) {
	c.routes = append(c.routes, route{method: "", pattern: pattern, handler: hdl})
}

func (c *ChiHandler) Build() (*http.Server, error) {
	srv := c.options.Server
	var router chi.Router = chi.NewRouter()
	router = router.With(c.options.Middlewares...)

	for _, r := range c.routes {
		switch r.method {
		case "":
			router.HandleFunc(r.pattern, r.handler)
		default:
			router.MethodFunc(r.method, r.pattern, r.handler)
		}
	}
	srv.Handler = router
	return srv, nil
}
