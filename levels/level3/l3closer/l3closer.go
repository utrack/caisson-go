package l3closer

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
	once    *sync.Once = &sync.Once{}
)

func Register(c Closer) {
	m.Lock()
	defer m.Unlock()
	closers = append(closers, c)
}

func Close(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	var retErr error
	once.Do(func() {
		for i := len(closers) - 1; i >= 0; i-- {
			errc := make(chan error, 1)

			go func() {
				errc <- closers[i].Close(ctx)
			}()
			select {
			case err := <-errc:
				if err != nil {
					retErr = err
				}
			case <-ctx.Done():
				retErr = ctx.Err()
			}
		}
	})
	return retErr
}
