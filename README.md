# caisson-go
Caisson is a meta-platform library for Go. A modular platform which you can pick up as is, or use it to build your own.

The layout of the repository is as follows:

- `pkg`: reusable, battle-tested blocks, which can be used to build your own platform libraries
- `errors`: a drop-in replacement, fully compatible with `errors` and `github.com/pkg/errors`; it extends the above libraries with primitives that are usually required by the server-side code