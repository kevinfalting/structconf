package structconf

import (
	"context"
	"errors"
)

// Default is a middleware which will set the default value for the field as
// long as no value was returned from the handler and the field is the zero
// value for it's type.
func Default() Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(
			func(ctx context.Context, f Field) (any, error) {
				result, err := h.Handle(ctx, f)
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
