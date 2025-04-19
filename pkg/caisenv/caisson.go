package caisenv

import (
	"context"

	_ "github.com/utrack/caisson-go/internal/caisenv"
	"github.com/utrack/caisson-go/internal/icloser"
	"github.com/utrack/caisson-go/log"
)

func Ensure() {
	// nothing to do for now, autoconfig
}

// Stop gracefully stops the environment, including anything registered via [github.com/utrack/caisson-go/closer].Register*.
//
// Please note that the closers are closed in LIFO order.
func Stop(ctx context.Context) error {
	log.Warn(ctx, "caisenv.Stop() called - stopping the environment", "module", "utrack/caisson-go")
	return icloser.Close(ctx)
}
