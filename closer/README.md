# caisson-go/closer

This is a platform-standard package to register things that should be closed before the app's shutdown.

It should be used for flushing/stopping any external connections (DBs, Kafka, custom metric flushes, etc).

Also, it can be used to close any internal loops gracefully.

## Order of the closers

Closers are closed in LIFO order - so in this example `bizlogic` will close first, allowing it to flush any info to the database before `db` closes:

```go
closer.RegisterC(db)
// ...
closer.RegisterC(bizlogic)
```

### Timeouts and contexts

If you register a `CloserC`, then the closing function will receive a context open for at most `CAISSON_GRACE_SHUTDOWN_TIMEOUT`(envvar) seconds - 30 seconds by default.

`CAISSON_GRACE_SHUTDOWN_TIMEOUT` is a total timeout for all closers' graceful shutdown - so, if the first closer takes 10 seconds to close, the second one will have only (30-10) seconds before it is killed.
