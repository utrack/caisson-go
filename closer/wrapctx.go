package closer

import "context"

type closeWrapper struct {
	f func() error
}

func (c *closeWrapper) Close() error {
	return c.f()
}

type closeCtxWrapper struct {
	f func(context.Context) error
}

func (c *closeCtxWrapper) Close(ctx context.Context) error {
	return c.f(ctx)
}

func RegisterFunc(f func() error) {
	Register(&closeWrapper{f: f})
}

func RegisterFuncC(f func(context.Context) error) {
	RegisterC(&closeCtxWrapper{f: f})
}
