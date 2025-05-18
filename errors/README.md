# caisson-go/errors

Package errors is a wrapper around the stdlib errors and pkg/errors.

It provides a consistent interface for error handling and enriching errors with details like HTTP codes, user messages, error contexts etc.

The intended usage is to use this package as a drop-in replacement for the stdlib errors package. Use the linter to block imports of stdlib or other errors.

## Usage

Fit to be used by: general API servers.  
You probably want to use the [errorbag](https://pkg.go.dev/github.com/utrack/caisson-go/pkg/errorbag) instead if the [Coder](https://pkg.go.dev/github.com/utrack/caisson-go/errors#Coder) is too big/too small for your needs, or if you don't have an HTTP server.

See the [godoc page](https://pkg.go.dev/github.com/utrack/caisson-go/errors) for details and examples.

## Rationale

The stdlib errors package provides the basic error handling functionality. It is a simple interface with a single method, Error(), and a few helper functions like Is(), As(), and others.

The github.com/pkg/errors package provides a more feature-rich error handling functionality, with stack traces, cause unwrapping, etc.

However, both of them lack the ability to enrich errors with signalling metadata, like the desired HTTP/gRPC response code, user-viewable messages, etc.  
This package addresses this gap by providing a consistent interface for error handling and enriching errors with signalling metadata.

By having this package as a drop-in, you'll promote the usage of `Coder` or a similar primitive, which will remove the friction of error decoration for the developers.  
Which usually leads to more consistent errors :)

This package also provides the `CoderDetailer` primitive, which allows you to enrich errors with typed details.  
The `CoderDetailer[T]` API will ensure that **every error is detailed** with a `T` value, as well as **provide an API for your users** to extract the details.  
No need to guess if it's a `pgconn.PgError` or `*pgconn.PgError` or `*sql.ErrNoRows` or `*MyCustomError`!

## Migration

Replace the `errors` or `github.com/pkg/errors` imports with `github.com/utrack/caisson-go/errors`; your code should continue to work as before.  
Start decorating your errors with metadata:

```go
// instead of
var ErrNotFound = errors.New("not found")
// use
var ErrNotFound = errors.NewCoder("NOT_FOUND").WithHTTPCode(404).WithMessage("Not found (user message)")

// ...

// in the function:
err := sql.ErrNoRows

return ErrNotFound.Wrap(err)
```

To retrieve the metadata, use the `Code` function:

```go

code := errors.Code(err)
if code != nil {
    //...
}

```
