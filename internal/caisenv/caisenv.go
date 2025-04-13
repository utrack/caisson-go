package caisenv

import "log/slog"

func init() {
	initTracer()
	slog.Info("caisenv initialized")
}