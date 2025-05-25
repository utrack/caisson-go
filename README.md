# caisson-go [![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/utrack/caisson-go)
Caisson is a meta-platform library for Go. A modular platform which you can pick up as is, or use it to build your own.

The layout of the repository is as follows:

- [caiapp](https://github.com/utrack/caisson-go/blob/main/caiapp/): a reference platform. To be split to a separate package.
- [closer](https://github.com/utrack/caisson-go/blob/main/closer/): a global closer registry, which can be used to close/flush resources on exit. 
- [errors](https://github.com/utrack/caisson-go/blob/main/errors/): a drop-in replacement, fully compatible with `errors` and `github.com/pkg/errors`. 
- [log](https://github.com/utrack/caisson-go/blob/main/log/): a structured logging package; enforces the usage of context for logging. 
- `pkg`: reusable, battle-tested blocks, which can be used to build your own platform libraries