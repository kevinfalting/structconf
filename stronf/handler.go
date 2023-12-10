package stronf

import (
	"context"
)

// HandlerFunc is a function type that implements the [Handler] interface. It
// takes a [context.Context], a [Field], and a interimValue as arguments. The
// returned value is the new interimValue for the [Field], or an error. For more
// information on how interimValue is used, see the the [Handler] interface.
type HandlerFunc func(ctx context.Context, field Field, interimValue any) (any, error)

// Handle allows [HandlerFunc] to implement the [Handler] interface. This method
// executes the [HandlerFunc] with the provided [context.Context], [Field], and
// interimValue.
func (h HandlerFunc) Handle(ctx context.Context, field Field, interimValue any) (any, error) {
	return h(ctx, field, interimValue)
}

var _ Handler = (HandlerFunc)(nil)

// Handler represents a contract for processing a [Field]'s proposed value in the
// context of configuration.
type Handler interface {
	// Handle method accepts a context, a Field, and an interimValue. It returns
	// an updated interimValue, or an error. The handler is responsible for
	// returning a new interimValue or the one it recieved.
	Handle(context.Context, Field, any) (any, error)
}

// CombineHandlers will return a handler that calls each handler in the order
// they are provided. The last handler takes precedence over the one before it.
func CombineHandlers(handlers ...Handler) Handler {
	var combinedHandlerFunc HandlerFunc = func(ctx context.Context, field Field, interimValue any) (any, error) {
		for _, h := range handlers {
			result, err := h.Handle(ctx, field, interimValue)
			if err != nil {
				return nil, err
			}

			interimValue = result
		}

		return interimValue, nil
	}

	return combinedHandlerFunc
}
