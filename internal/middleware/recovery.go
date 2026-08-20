package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/MdRasB/LogLine/internal/contextutil"
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
					requestID, _ := contextutil.RequestID(r.Context())

					logger.Error(
						"panic recovered",
						"request_id", requestID,
						"method", r.Method,
						"path", r.URL.Path,
						"panic", err,
						"stack", string(debug.Stack()),
					)

					WriteJSON(w, http.StatusInternalServerError, map[string]string{
						"error": "internal server error",
					})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
