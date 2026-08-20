// Package middleware handle the middleware logics for the routes and servers
package middleware

import (
	"net/http"

	"github.com/MdRasB/LogLine/internal/auth"
	"github.com/MdRasB/LogLine/internal/contextutil"
)

func AuthMiddleware(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			sessiontoken := auth.ExtractBearerToken(authHeader)

			if sessiontoken == "" {
				WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "missing session token",
				})
				return
			}

			session, err := authService.ValidateSession(sessiontoken)
			if err != nil {
				WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid session token",
				})
				return
			}

			ctx := contextutil.WithUserID(
				r.Context(),
				session.UserID,
			)

			ctx = contextutil.WithSessionID(
				ctx,
				session.ID,
			)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
