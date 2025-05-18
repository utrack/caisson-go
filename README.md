# caisson-go
Caisson is a meta-platform library for Go. A modular platform which you can pick up as is, or use it to build your own.

The layout of the repository is as follows:

- `caiapp`: a reference platform. To be split to a separate package.
- `closer`: a global closer registry, which can be used to close/flush resources on exit. [README](https://github.com/utrack/caisson-go/blob/main/closer/README.md)
- `errors`: a drop-in replacement, fully compatible with `errors` and `github.com/pkg/errors`. [README](https://github.com/utrack/caisson-go/blob/main/errors/README.md)
- `log`: a structured logging package; enforces the usage of context for logging. Compatible with `slog` and `zap`. [README](https://github.com/utrack/caisson-go/blob/main/log/README.md)
- `pkg`: reusable, battle-tested blocks, which can be used to build your own platform libraries