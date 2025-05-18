# caisson-go/log

Package log provides a structured logging package; enforces the usage of context for logging. Compatible with `slog` and `zap`.

## Usage

Suitable for general use. Minimally opinionated and extendable.  
Make sure to set up the global `slog` logger; example is at [internal/caisenv/caisenv.go](https://github.com/utrack/caisson-go/blob/main/internal/caisenv/caisenv.go#L19).

TODO godoc

## Rationale

The stdlib `log`/`slog` packages do not enforce the usage of context for logging, leading to nasty logs in production.

Such logging leads to nasty anti-patterns like this:
```go
// ExcessiveParameters embeds the userID in every log.
// The ctx is not passed so EVERY log will not contain the trace ID.
func ExcessiveParameters() error {
    slog.Info("user logged in", "userID",userID)
    // ... 
    slog.Error("failed to do something", "error",err,"userID",userID)
    // ...
    // which leads to these mistakes when the userID is forgotten:
    slog.Warn("doing something else")
    // or when the error is logged under some other key:
    slog.Error("failed to do something", "error",err)
}
```

Generally I've found these issues when I needed to firefight something, and the logs were not helpful enough.

This package proposes a) mandatory context usage and b) a consistent API for logging errors.

### Mandatory context usage

Every logging function requires a context, which is used to propagate the log context to downstream code.  
The underlying logger can extract anything from the context, not only keys propagated via `With()`.  
As an example, [slogtrace](https://github.com/utrack/caisson-go/blob/main/pkg/slogtrace/slogadapter.go) extracts the trace/span IDs from the context and adds them to the log.

The `log` API ensures that you won't lose anything!

Usage example:
```go
func main() error {
    ctx := context.Background()
    // every log will contain the 'module'
    return KafkaConsumerLoop(log.With(ctx, "module", "kafkaConsumer"))
}

func KafkaConsumerLoop(ctx context.Context) error {
    for {
        err := kafkaReadMessage(ctx)
        if err != nil {
            // logs {error.message: failed to read messages, error.stack: ..., module: kafkaConsumer}
            log.Error(ctx, "failed to read messages", err)
        }
    }
}

func kafkaReadMessage(ctx context.Context) error {
    // logs {info: doing something, module: kafkaConsumer}
    log.Info(ctx, "doing something")
    return nil
}

```

### Consistent API for logging errors

The `log` API ensures that the error is always logged with the same key and form.  
`log.Error()` has a mandatory `error` parameter, which is logged with the key `error.message` and `error.stack`.

You can still use `log.Errorn()` if you don't have an error to log, but you still want to log a message with the ERROR level.  
This is a good example of added friction - it's easy to do what's right (`log.Error` is intuitive), while `log.Errorn/Errorne` look foreign and require additional thought.