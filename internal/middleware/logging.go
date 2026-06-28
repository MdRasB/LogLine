package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			rw := NewResponseWriter(w)

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			requestID, _ := GetRequestID(r.Context())

			logger.Info(
				"http request",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration", duration.String(),
			)
		})
	}
}

