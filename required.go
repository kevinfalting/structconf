package structconf

import (
	"context"
	"fmt"
	"reflect"
)

// Required is a middleware that will a value to be set. Zero values are allowed
// when returned from handlers.
func Required() Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(
			func(ctx context.Context, f Field) (any, error) {
				result, err := h.Handle(ctx, f)
				if err != nil {
					return nil, err
				}

				if result != nil && isZero(result) {
					return result, nil
				}

				if !f.IsZero() {
					return result, nil
				}

				_, required := f.LookupTag("conf", "required")

				if required && result == nil {
					return nil, fmt.Errorf("required field %s is not set", f.Name())
				}

				return result, nil
			},
		)
	}
}

func isZero(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Func:
		return rv.IsNil()
	case reflect.Map, reflect.Slice:
		return rv.IsNil() || rv.Len() == 0
	case reflect.Array:
		z := true
		for i := 0; i < rv.Len(); i++ {
			z = z && isZero(rv.Index(i).Interface())
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < rv.NumField(); i++ {
			z = z && isZero(rv.Field(i).Interface())
		}
		return z
	}

	// Compare other types directly:
	z := reflect.Zero(rv.Type())
	result := rv.Interface() == z.Interface()

	return result
}
