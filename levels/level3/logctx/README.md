# logctx

This package provides a way to pass loggers down from the caller to the callee via the context.

It uses the stdlib `slog` package as a logger; swap out the default `slog` logger with your own implementation if you need to.

## Design and reasoning

## Intended usage

The caller-to-be `log` package should wrap the `With()` function and provide it to the application code.  
The application code should NOT use the `FromCtx()` function directly; instead, the wrapping `log` package should call it for every log line.

The `log` package's logging functions should also accept `ctx` as the first argument for every logging call.

This way, the application code won't need to know about the `log`'s intricacies, and the logs will be enriched with the context data as long as the `ctx` is passed.

## Rationale

Passing down the logger via the context lets you get rid from these patterns:

#### Passing the logger as an argument to the function

It is very noisy and too explicit for something like logging. People get complacent and start recreating the logger in every major module, which drops the context for the callee's logs.

#### Injecting the logger into the struct

Those implementations usually muddle the structs, adding a logger field to every struct that needs logging. This is a bad practice, as it pollutes the struct's API and makes it harder to reason about the code.

#### Using a global logger everywhere without context forced by the caller

This doesn't let you gradually enrich the log data, unless you pass a lot of the log-only data to the callee. 

## Examples

See the [level6/log](../../level6/log) package for an example of a higher-level logging package.