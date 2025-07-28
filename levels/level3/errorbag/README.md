# errorbag

This package provides a way to enrich errors with arbitrary keys and values.  
It is designed to be used as a building block for other packages, allowing them to add context to the errors.

Generally it is not intended to be used directly; instead, use it in your `errors` drop-in package to provide a standard way of error decoration.

## Reasoning

While the stdlib `context` package provides a way to pass values down from the caller to the callee, there's no standard way to do the same for bubbling up error context.

This package works almost exactly like `context`, but reversed - allowing you to decorate your error when it bubbles up from the callee to the caller.



