# l3closer

This package provides a global Closer registry. A Closer is something that should be closed when the application is shutting down.

l3closer's implementation is bare bones; it does not provide any concurrency safety, nor does it provide any coordination between closing and the graceful shutdown.
Those features must be provided by the higher Levels, or the custom platform code. The higher level code must call `l3closer.Close(ctx)` when the application is shutting down.

The registered closers are closed in the LIFO order, allowing you to follow the dependency order.