package icloser

import (
	"context"
	"sync"
)

type Closer interface {
	Close(ctx context.Context) error
}

var (
	closers []Closer
	m       sync.Mutex
)

func Register(c Closer) {
	m.Lock()
	defer m.Unlock()
	closers = append(closers, c)
}

func Close(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	for i := len(closers) - 1; i >= 0; i-- {
		errc := make(chan error, 1)

		go func() {
			errc <- closers[i].Close(ctx)
		}()
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
