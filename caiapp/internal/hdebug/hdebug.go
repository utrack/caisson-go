/*
Package hdebug provides a separate server
which serves the livez/readyz/healthz endpoints,
as well as pprof handlers and Prometheus metrics.
*/
package hdebug

import (
	"expvar"
	"html/template"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"path"
	"strings"

	"github.com/felixge/fgprof"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/utrack/caisson-go/caiapp/handler"
	"github.com/utrack/caisson-go/caiapp/internal/hchi"
	"github.com/utrack/caisson-go/pkg/http/hhandler"
)

type Mux struct {
	mux   hhandler.Server
	ready *atomBool
}

func (m *Mux) Build() (*http.Server, error) {
	return m.mux.Build()
}

func (m *Mux) SetReady(ready bool) {
	m.ready.Set(ready)
}

// New returns an HTTP server for internal usage.
// It serves profiling info, docs and Prometheus metrics.
func New() *Mux {
	mux := hchi.New()

	atom := &atomBool{}

	// Disable caching for all responses
	// Nothing here should be cached
	mux.Apply(handler.WithGlobalMiddleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				next.ServeHTTP(w, r)
			})
		},
	))

	mux.MethodFunc("GET", "/metrics", promhttp.Handler().ServeHTTP)

	mux.HandleFunc("/debug/fgprof", fgprof.Handler().ServeHTTP)
	mux.HandleFunc("/debug/vars", expvar.Handler().ServeHTTP)
	mux.HandleFunc("/debug/bin", func(w http.ResponseWriter, r *http.Request) {
		pp, err := os.Executable()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fstream, err := os.OpenFile(pp, os.O_RDONLY, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fstream.Close()
		filename := strings.TrimPrefix(path.Base(pp), pp)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		_, err = io.Copy(w, fstream)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	mux.HandleFunc("/debug/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/debug/pprof/*", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`"ok"`))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if atom.Get() {
			_, _ = w.Write([]byte(`"ok"`))
			return
		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// TODO _ = json.NewEncoder(w).Encode(ai)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		err := tplDebugHome.Execute(w, struct {
			//TODO App appinfo.Info
		}{
			//App: ai,
		})
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
	})

	// TODO serve docs/swagger
	return &Mux{mux: mux, ready: atom}
}

var tplDebugHome = template.Must(template.New("foo").Parse(
	`
<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title>debug port for {{ .App.Name }}</title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link rel="stylesheet" href="https://unpkg.com/chota@0.8.1/dist/chota.min.css">
        <!-- Place favicon.ico in the root directory -->

    </head>
    <body>
<div class="container">        
<h1>debug port for {{ .App.Name }}</h1>
<blockquote>
<p>
<b>name</b>: {{ .App.Name }}<br>
<b>version</b>: {{ .App.Version }}<br>
<b>build date</b>: {{ .App.BuildDate }}<br>
<b>git log</b>: {{ .App.GitLog }}<br>
<b>go version</b>: {{ .App.GoVersion }}<br>
</p>
</blockquote>
<h2>Available handlers</h2>
<ul>
<li>/docs - go to (primary HTTP server)/docs to see the Swagger UI</li>
<li>/redoc - go to (primary HTTP server)/redoc to see the alternative Redoc API documentation</li>
<li><a href="/grpcui/">/grpcui</a> - gRPC UI (may be unavailable if no gRPC service is enabled)</li>
<li><a href="/version">/version</a> - version info in JSON format</li>
</ul>
<ul>
<li><a href="/debug/bin">/debug/bin</a> - This build's executable file</li>
<li><a href="/debug/pprof">/debug/pprof</a> - Go pprof</li>
<li><a href="/debug/fgprof">/debug/fgprof</a> - github.com/felixge/fgprof prof dump</li>
<li><a href="/debug/vars">/debug/vars</a> - Go expvar</li>
</ul>
</div>
    </body>
</html>
`))
