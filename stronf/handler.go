package stronf

import (
	"context"
)

// HandlerFunc is a function type that implements the [Handler] interface. It
// takes a [context.Context], a [Field], and a proposedValue as arguments. The
// returned value is the new proposedValue for the [Field], or an error. For more
// information on how proposedValue is used, see the the [Handler] interface.
type HandlerFunc func(ctx context.Context, field Field, proposedValue any) (any, error)

// Handle allows [HandlerFunc] to implement the [Handler] interface. This method
// executes the [HandlerFunc] with the provided [context.Context], [Field], and
// proposedValue.
func (h HandlerFunc) Handle(ctx context.Context, field Field, proposedValue any) (any, error) {
	return h(ctx, field, proposedValue)
}

var _ Handler = (HandlerFunc)(nil)

// Handler represents a contract for processing a [Field]'s proposed value in the
// context of configuration.
type Handler interface {
	// Handle method accepts a context, a Field, and an proposedValue. It returns
	// an updated proposedValue, or an error. The handler is responsible for
	// returning a new proposedValue or the one it recieved.
	Handle(context.Context, Field, any) (any, error)
}

// CombineHandlers will return a handler that calls each handler in the order
// they are provided. The last handler takes precedence over the one before it.
func CombineHandlers(handlers ...Handler) Handler {
	var combinedHandlerFunc HandlerFunc = func(ctx context.Context, field Field, proposedValue any) (any, error) {
		for _, h := range handlers {
			result, err := h.Handle(ctx, field, proposedValue)
			if err != nil {
				return nil, err
			}

			proposedValue = result
		}

		return proposedValue, nil
	}

	return combinedHandlerFunc
}
