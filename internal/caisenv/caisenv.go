package caisenv

import (
	"log/slog"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
)

func init() {
	olog := logr.FromSlogHandler(slog.Default().Handler())
	otel.SetLogger(olog)
	initTracer()
	initMetrics()
	slog.Info("caisenv initialized")
}
