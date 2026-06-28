package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	//"time"
)

func Recovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			defer func() {
				if err := recover(); err != nil {
					requestID, _ := GetRequestID(r.Context())

					logger.Error(
						"panic recovered",
						"request_id", requestID,
						"method", r.Method,
						"path", r.URL.Path,
						"panic", err,
						"stack", string(debug.Stack()),
					)

					http.Error(
						w, "internel server error", http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

