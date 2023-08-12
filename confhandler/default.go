package confhandler

import (
	"context"
	"errors"

	"github.com/kevinfalting/structconf/stronf"
)

// Default is a middleware which will set the default value for the field as
// long as no value was returned from the handler and the field is the zero
// value for it's type.
func Default() stronf.Middleware {
	return func(h stronf.Handler) stronf.Handler {
		return stronf.HandlerFunc(
			func(ctx context.Context, f stronf.Field, interimValue any) (any, error) {
				result, err := h.Handle(ctx, f, interimValue)
				if err != nil {
					return nil, err
				}

				if result != nil {
					return result, nil
				}

				if !f.IsZero() {
					return result, nil
				}

				defaultVal, ok := f.LookupTag("conf", "default")
				if !ok {
					return result, nil
				}

				if len(defaultVal) == 0 {
					return nil, errors.New("empty default value")
				}

				return parseStringForKind(defaultVal, f.Kind())
			},
		)
	}
}
