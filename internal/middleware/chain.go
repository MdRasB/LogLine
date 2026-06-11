package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(
	h http.Handler,
	middlewares ...Middleware,
) http.Handler {
	wrapped := h

	for i := len(middlewares)-1; i>=0 ; i-- {
		wrapped = middlewares[i](wrapped)
	}

	return wrapped
}
