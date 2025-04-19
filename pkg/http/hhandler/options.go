package hhandler

import (
	"net/http"

	"crypto/tls"
)

// Options struct is intended to be used internally by the framework.
//
// The framework itself needs to set the defaults.
// The framework may implement its own custom options and/or alias the OptionHTTP type(s).
type Options struct {
	Server *http.Server

	TlsConfig   *tls.Config
	Middlewares []func(http.Handler) http.Handler
}

type OptionHTTP func(o *Options)
