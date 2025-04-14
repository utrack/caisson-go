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
	slog.SetDefault(slog.New(slogtrace.NewContextHandler(slog.NewJSONHandler(os.Stdout, nil))))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	olog := logr.FromSlogHandler(slog.Default().Handler())
	otel.SetLogger(olog)
	initTracer()
	initMetrics()

	slog.Info("caisson-go environment initialized", "config", caisconfig.Get())
}
