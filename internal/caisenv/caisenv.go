package caisenv

import (
	"log/slog"

	"github.com/go-logr/logr"
	"github.com/utrack/caisson-go/pkg/caisconfig"
	"go.opentelemetry.io/otel"
)

func init() {
	olog := logr.FromSlogHandler(slog.Default().Handler())
	otel.SetLogger(olog)
	initTracer()
	initMetrics()

	slog.Info("caisson-go environment initialized", "config", caisconfig.Get())
}
