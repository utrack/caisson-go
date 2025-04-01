package hserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/plog"
	"golang.org/x/sync/errgroup"
)

// Server is an HTTP server that's already accepting the connections.
//
// It is designed to be ingress-friendly and used in a container environment, k8s or otherwise.
//
// Use Ready() as an another healthcheck for the ingress.
//
// As soon as it is created, it starts accepting the connections. The connections hang in a waiting state until the Run() call; then it passes the requests to the provided HTTP server.
//
// Keep in mind that there's no delay between Server.GracefulShutdown() and the server being stopped.
// To ensure zero downtime for your infra, mark some other healthcheck as not ready and wait some time (healthcheck period*2) before calling Server.GracefulShutdown().
//
// This will ensure that the ingress will stop routing traffic to the server before it stops accepting the connections.
type Server struct {
	lis  net.Listener
	opts opts

	readyToServe *atomBool

	eg    *errgroup.Group
	egCtx context.Context

	srv *http.Server
}

type opts struct {
	// addr is a net address to listen on
	// (e.g. "localhost:8080", ":8080", "127.0.0.1:8080", etc)
	addr string

	name string
}

type Option func(*opts)

// WithName sets a name for the server that'll appear in the logs.
//
// Defaults to "primary".
func WithName(name string) Option {
	return func(o *opts) {
		o.name = name
	}
}

// New creates a new Server and starts listening immediately.
// It is non-blocking; New() returns as soon as the listener is established.
// Returns an error if the listener cannot be established.
//
// addr is a net address to listen on (e.g. "localhost:8080", ":8080", "127.0.0.1:8080", etc)
func New(addr string, optionFuncs ...Option) (*Server, error) {
	o := &opts{
		addr: addr,
		name: "primary",
	}
	for _, opt := range optionFuncs {
		opt(o)
	}
	lis, err := net.Listen("tcp", o.addr)
	if err != nil {
		return nil, errors.Wrapd(err, "failed to bring up the listener", "addr", o.addr, "name", o.name)
	}
	plog.Info(context.Background(), "listener port bound", "addr", o.addr, "name", o.name)

	eg, egCtx := errgroup.WithContext(context.Background())

	return &Server{
		lis:          lis,
		opts:         *o,
		eg:           eg,
		egCtx:        egCtx,
		readyToServe: &atomBool{},
	}, nil
}

// Run starts serving the HTTP requests.
// Blocks until the server stops or the context is canceled.
// If a context is canceled, the server will be stopped immediately.
func (s *Server) Run(ctx context.Context, h *http.Server) error {

	// start serving requests
	s.eg.Go(func() error {
		if h.TLSConfig != nil {
			return h.ServeTLS(s.lis, "", "")
		}
		return h.Serve(s.lis)
	})

	// stop the server when the context passed to Run is canceled
	s.eg.Go(func() error {
		select {
		// happens when any of the other goroutines returns an error
		case <-s.egCtx.Done():
			return nil
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "Run() context is canceled")
		}
	})

	// mark the server as ready after a second.
	// this is to ensure that h.Serve() call ran successfully.
	s.eg.Go(func() error {
		<-time.After(time.Second)
		// if an error happened, don't signal that the server is ready
		if s.egCtx.Err() != nil {
			return nil
		}

		plog.Info(ctx, "HTTP server is ready and accepting connections", "name", s.opts.name)

		s.readyToServe.Set(true)

		return nil
	})

	err := s.eg.Wait()

	// if an error happened, immediately set readyToServe to false
	s.readyToServe.Set(false)

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Ready() bool {
	return s.readyToServe.Get()
}

// GracefulStop stops the HTTP server gracefully.
// Blocks until the server stops or the context is canceled.
// If a context is canceled, the server will be stopped immediately.
func (s *Server) GracefulStop(ctx context.Context) error {

	s.readyToServe.Set(false)

	s.eg.Go(func() error {
		return errors.Wrap(s.srv.Shutdown(ctx), "failed to gracefully stop the server")
	})

	err := s.eg.Wait()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
