package caisenv

import (
	"context"

	"github.com/utrack/caisson-go/internal/caisenv"
)

func Ensure() {
	// nothing to do for now, autoconfig
}

func init() {
	caisenv.Ensure()
}

// Stop gracefully stops the environment, including anything registered via [github.com/utrack/caisson-go/closer].Register*.
//
// Please note that the closers are closed in LIFO order.
func Stop(ctx context.Context) error {
	return caisenv.Stop(ctx)
}
