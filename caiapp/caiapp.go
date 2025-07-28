package caiapp

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	otelchimetric "github.com/riandyrn/otelchi/metric"
	"github.com/utrack/caisson-go/caiapp/internal/cappconfig"
	"github.com/utrack/caisson-go/caiapp/internal/hchi"
	"github.com/utrack/caisson-go/caiapp/internal/hdebug"
	"github.com/utrack/caisson-go/caiapp/internal/oapigen"
	"github.com/utrack/caisson-go/caiapp/internal/sdescbind"
	"github.com/utrack/caisson-go/closer"
	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/levels/level3/servers/l3http"
	"github.com/utrack/caisson-go/log"
	"github.com/utrack/caisson-go/pkg/caisenv"
	"github.com/utrack/caisson-go/pkg/http/hhandler"
	"github.com/utrack/caisson-go/pkg/plconfig"
	"github.com/utrack/pontoon/sdesc"
	"golang.org/x/sync/errgroup"
)

type App struct {
	handlers *Handlers
	hsrv     *l3http.Server
	setReady func(bool)
	eg       *errgroup.Group
	egCtx    context.Context
}

func New() (*App, error) {

	caisenv.Ensure()

	cfg, err := cappconfig.Get()
	if err != nil {
		return nil, errors.Wrap(err, "when configuring caiapp")
	}

	debugListener, err := l3http.New(cfg.Server.AddrDebug+":"+strconv.Itoa(cfg.Server.PortDebug), l3http.WithName("debug"))
	if err != nil {
		return nil, errors.Wrap(err, "when creating a debug HTTP server")
	}

	debugMux := hdebug.New()
	debugHandler, err := debugMux.Build()
	if err != nil {
		return nil, errors.Wrap(err, "when creating a debug HTTP handler")
	}

	eg, egCtx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		// use Background since we don't want to stop debug ever
		return debugListener.Run(context.Background(), debugHandler)
	})

	mainSrv, err := l3http.New(cfg.Server.AddrHTTP+":"+strconv.Itoa(cfg.Server.PortHTTP), l3http.WithName("main"))
	if err != nil {
		return nil, errors.Wrap(err, "when creating main HTTP server")
	}

	mainHdl := hchi.New()

	return &App{
		hsrv:     mainSrv,
		handlers: &Handlers{http: mainHdl},
		setReady: debugMux.SetReady,
		eg:       eg,
		egCtx:    egCtx,
	}, nil
}

func (a *App) Handlers() *Handlers {
	return a.handlers
}

func (a *App) Run(ctx context.Context, services ...sdesc.Service) error {

	cfg, err := cappconfig.Get()
	if err != nil {
		return errors.Wrap(err, "when reading caisson caiapp config")
	}

	ctx = log.With(ctx, "module", "caiapp")

	log.Info(ctx, "caiapp starting")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	caiconf := plconfig.Get()

	otelChiCfg := otelchimetric.NewBaseConfig(caiconf.ServiceName)
	hsrv := a.handlers.http
	// Prepend critical middlewares - request logger,
	// metrics, recovery, tracer etc.
	// They should go in front of the app-provided middlewares.
	hsrv.Apply(func(o *hhandler.Options) {
		o.Middlewares = append([]func(http.Handler) http.Handler{

			chimw.RealIP,
			otelchi.Middleware(caiconf.ServiceName,
				otelchi.WithTraceResponseHeaders(otelchi.TraceHeaderConfig{
					TraceIDHeader:      "X-Trace-Id",
					TraceSampledHeader: "X-Trace-Sampled",
				}),
			),
			otelchimetric.NewRequestDurationMillis(otelChiCfg),
			otelchimetric.NewRequestInFlight(otelChiCfg),
			otelchimetric.NewResponseSizeBytes(otelChiCfg),
			chimw.Recoverer,
		}, o.Middlewares...)
	})

	handlerDocMeta := []oapigen.HandlerDesc{}
	for i, s := range services {
		hdl, err := sdescbind.Bind(s, a.handlers.http)
		if err != nil {
			return errors.Wrapf(err, "when binding HTTP handlers for service %d (%T)", i, s)
		}
		handlerDocMeta = append(handlerDocMeta, hdl...)
	}

	hExts := hsrv.Extensions()
	doc, err := oapigen.GenerateOAPI(handlerDocMeta, hExts)
	if err != nil {
		return errors.Wrap(err, "when generating OpenAPI document")
	}

	// TODO move away
	hsrv.MethodFunc("GET", "/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/openapi+yaml")
		b, err := doc.Render()
		if err != nil {
			log.Error(r.Context(), "when rendering OpenAPI document", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
	})

	oapiPath := path.Join(hExts.Prefix, "/openapi.yaml")

	hsrv.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		ret := `
		<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>Elements in HTML</title>
  
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
  </head>
  <body>

    <elements-api
      apiDescriptionUrl="` + oapiPath + `"
      router="hash"
    />

  </body>
</html>
		
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(ret))

	})

	finalHandler, err := hsrv.Build()
	if err != nil {
		return errors.Wrap(err, "when building main HTTP handler")
	}

	closer.RegisterFuncC(a.hsrv.GracefulStop)

	a.eg.Go(func() error {
		err := a.hsrv.Run(ctx, finalHandler)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return errors.Wrap(err, "when running main HTTP server")
	})

	// wait for possible errors from the handler
	select {
	case <-time.After(100 * time.Millisecond):
	case <-a.egCtx.Done():
		return errors.Wrap(a.egCtx.Err(), "when running main HTTP server")
	}

	a.setReady(true)

	select {
	case <-a.egCtx.Done():
	case <-ctx.Done():
		log.Info(ctx, "caiapp: app.Run() context canceled", "reason", ctx.Err())
	case sig := <-sigs:
		log.Info(ctx, "caiapp: caught signal", "signal", sig)
	}

	log.Warn(ctx, "initiating graceful shutdown, Ctrl-C again to force exit", "grace_delay", cfg.GracefulShutdown.Delay, "grace_timeout", cfg.GracefulShutdown.Timeout)

	// die on ^C^C
	go func() {
		sig := <-sigs
		log.Fatal(ctx, "second signal,terminating", "sig", sig.String())
	}()
	a.setReady(false)

	<-time.After(cfg.GracefulShutdown.Delay)
	log.Info(ctx, "graceful shutdown delay expired, shutting down")

	return caisenv.Stop(ctx)
}
