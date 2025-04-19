package hhandler

import (
	"net/http"
)

// Configurer is intended to be used by the application code.
//
// It sets up global server options.
//
// It does not expose the internal server configuration so that the framework will control it and its defaults.
type Configurer interface {
	Apply(oo ...OptionHTTP)
}

// Server is a dropped-priveleges Handler that only creates a final http.Server.
//
// It lets you delay the server creation until all the global options are applied by the application code via [Configurer].
type Server interface {
	Handler
	Configurer
	Build() (*http.Server, error)
}

// Handler should be used internally by the framework
// and not by the application code.
//
// It allows the framework to set up all the application
// routes and handlers.
type Handler interface {
	MethodFunc(method string, pattern string, hdl http.HandlerFunc)
	HandleFunc(pattern string, hdl http.HandlerFunc)
}
