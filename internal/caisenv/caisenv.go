package caisenv

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-logr/logr"
	"github.com/utrack/caisson-go/closer"
	"github.com/utrack/caisson-go/internal/icloser"
	"github.com/utrack/caisson-go/log"
	"github.com/utrack/caisson-go/pkg/plconfig"
	"github.com/utrack/caisson-go/pkg/slogtrace"
	"go.opentelemetry.io/otel"
)

func init() {

	// slogtrace extracts trace_id/span_id from the context. Use it for the global logger.
	handler := slogtrace.NewContextHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	logger := slog.New(handler)
	slog.SetDefault(logger)

	olog := logr.FromSlogHandler(handler)
	otel.SetLogger(olog)
	closeTracer := initTracer()
	closeMetrics := initMetrics()

	closer.RegisterFuncC(closeTracer)
	closer.RegisterFuncC(closeMetrics)

	slog.Info("caisson-go environment initialized", "config", plconfig.Get())
}

// Stop gracefully stops the environment, including anything registered via [github.com/utrack/caisson-go/closer].Register*.
//
// Please note that the closers are closed in LIFO order.
func Stop(ctx context.Context) error {
	log.Warn(ctx, "caisenv.Stop() called - stopping the environment", "module", "utrack/caisson-go")
	return icloser.Close(ctx)
}
