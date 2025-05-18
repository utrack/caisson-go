/*
Package errors is a wrapper around the stdlib errors and pkg/errors.

It provides a consistent interface for error handling and enriching errors with details like HTTP codes, user messages, error contexts etc.

The intended usage is to use this package as a drop-in replacement for the stdlib errors package. Use the linter to block imports of stdlib or other errors.
*/
package errors
