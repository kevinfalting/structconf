package structconf

import "context"

type Middleware func(h Handler) Handler

func WrapMiddleware(handlers []Handler, mw ...Middleware) Handler {
	var handler Handler = HandlerFunc(
		func(ctx context.Context, f Field) (any, error) {
			var val any

			for _, h := range handlers {
				result, err := h.Handle(ctx, f)
				if err != nil {
					return nil, err
				}

				if result == nil {
					continue
				}

				val = result
			}

			return val, nil
		},
	)

	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
