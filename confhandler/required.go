package confhandler

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kevinfalting/structconf/stronf"
)

// Required is a middleware that will require a value to be set. Zero values are
// allowed when returned from handlers.
func Required() stronf.Middleware {
	return func(h stronf.Handler) stronf.Handler {
		return stronf.HandlerFunc(
			func(ctx context.Context, f stronf.Field, interimValue any) (any, error) {
				result, err := h.Handle(ctx, f, interimValue)
				if err != nil {
					return nil, err
				}

				if result != nil && reflect.ValueOf(result).IsZero() {
					return result, nil
				}

				if !f.IsZero() {
					return result, nil
				}

				_, required := f.LookupTag("conf", "required")

				if required && result == nil {
					return nil, fmt.Errorf("structconf: required field %s is not set", f.Name())
				}

				return result, nil
			},
		)
	}
}
