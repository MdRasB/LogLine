package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/MdRasB/LogLine/internal/contextutil"
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

			requestID, _ := contextutil.RequestID(r.Context())
			userID, hasUserID := contextutil.UserID(r.Context())
			sessionID, hasSessionID := contextutil.SessionID(r.Context())

			args := []any{
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration", duration.String(),
			}

			if hasUserID {
				args = append(args, "user_id", userID)
			}

			if hasSessionID {
				args = append(args, "session_id", sessionID)
			}

			logger.Info(
				"http request",
				args...,
			)
		})
	}
}
