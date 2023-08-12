package stronf

import (
	"context"
)

// HandlerFunc is a function type that implements the Handler interface. It
// takes a context, a Field, and a interimValue as arguments. The returned value
// is the new interimValue for the Field, or an error. For more information on
// how interimValue is used, see the documentation for the Handler interface.
type HandlerFunc func(context.Context, Field, any) (any, error)

// Handle allows HandlerFunc to implement the Handler interface. This method
// executes the HandlerFunc with the provided context, Field, and interimValue.
func (h HandlerFunc) Handle(ctx context.Context, f Field, interimValue any) (any, error) {
	return h(ctx, f, interimValue)
}

var _ Handler = (HandlerFunc)(nil)

// Handler represents a contract for processing a Field's proposed value in the
// context of configuration.
type Handler interface {
	// Handle method accepts a context, a Field, and an interimValue. It returns
	// an updated interimValue, or an error. If a Handler wishes to suggest a
	// different interimValue to be considered by the next handler, it should
	// return a non-nil value. Returning nil indicates that the Handler does
	// not wish to modify the current interimValue.
	Handle(context.Context, Field, any) (any, error)
}
