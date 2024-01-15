package stronf

import (
	"context"
)

// HandleFunc is the function signature for working on a [Field].
type HandleFunc func(ctx context.Context, field Field, proposedValue any) (any, error)

// CombineHandlers will return a handler that calls each handler in the order
// they are provided. The last handler takes precedence over the one before it.
func CombineHandlers(handlers ...HandleFunc) HandleFunc {
	var combinedHandlerFunc HandleFunc = func(ctx context.Context, field Field, proposedValue any) (any, error) {
		for _, h := range handlers {
			result, err := h(ctx, field, proposedValue)
			if err != nil {
				return nil, err
			}

			proposedValue = result
		}

		return proposedValue, nil
	}

	return combinedHandlerFunc
}
