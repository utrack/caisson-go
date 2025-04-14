package caisenv

import (
	"log/slog"
	"os"

	"github.com/go-logr/logr"
	"github.com/utrack/caisson-go/pkg/caisconfig"
	"github.com/utrack/caisson-go/pkg/slogtrace"
	"go.opentelemetry.io/otel"
)

func init() {
	handler := slogtrace.NewContextHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	logger := slog.New(handler)
	slog.SetDefault(logger)

	olog := logr.FromSlogHandler(handler)
	otel.SetLogger(olog)
	initTracer()
	initMetrics()

	slog.Info("caisson-go environment initialized", "config", caisconfig.Get())
}
