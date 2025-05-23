package web

// Middleware represents an HTTP handler wrapped in another handler
type Middleware func(Handler) Handler

// wrapMiddleware creates a Handler by wrapping middleware around a final handler
// Ordered by how they come in the slice
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := mw[i]
		if mwFunc != nil {
			handler = mwFunc(handler)
		}
	}

	return handler
}
