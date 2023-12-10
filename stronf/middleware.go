package stronf

// Middleware is a function that takes a Handler and returns a modified or
// wrapped handler.
type Middleware func(handler Handler) Handler

// WrapMiddleware returns a handler that wraps the provided handler with the
// provided middleware. Middleware are applied in reverse order to achieve a
// call order that will call the first, second, third, on the way in, and third,
// second, first, on the way out.
func WrapMiddleware(handler Handler, middleware ...Middleware) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		mw := middleware[i]
		if mw == nil {
			continue
		}

		handler = mw(handler)
	}

	return handler
}
