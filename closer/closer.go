package closer

import (
	"context"
	"reflect"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/levels/level3/l3closer"
)

type Closer interface {
	Close() error
}

type CloserC interface {
	Close(ctx context.Context) error
}

type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

func Register(c Closer) {
	l3closer.Register(closeCtxWrap{c: c})
}

func RegisterC(c CloserC) {
	l3closer.Register(c)
}

type closeCtxWrap struct {
	c Closer
}

func (c closeCtxWrap) Close(ctx context.Context) error {
	eret := make(chan error, 1)

	go func() {
		eret <- c.c.Close()
	}()

	tn := reflect.TypeOf(c.c).Name()

	select {
	case err := <-eret:
		return err
	case <-ctx.Done():
		return errors.Wrapf(ctx.Err(), "global Closer couldn't close '%v' in time", tn)
	}
}

// ClosingContext creates a child Context which will close
// when the global Closer is being closed.
//
// Please note that usage of a closing context is discouraged
// since it does not wait for the underlying biz operation to finish --
// so it may happen that your business logic will work until after the
// DB is closed for example -- even if you've put the DB closer
// before ClosingContext.
func ClosingContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	cc := &closingContext{
		cancel: cancel,
	}

	Register(cc)
	return ctx
}

type closingContext struct {
	cancel func()
}

func (c *closingContext) Close() error {
	c.cancel()
	return nil
}
