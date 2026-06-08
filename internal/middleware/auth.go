package middleware

import (
	"context"
	"net/http"

	"github.com/MdRasB/LogLine/internal/auth"
)

type contextKey string

var (
	UserIDKey    = contextKey("user_id")
	SessionIDKey = contextKey("session_id")
)

func AuthMiddleware(
	authService *auth.Service,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			sessiontoken := auth.ExtractBearerToken(authHeader)

			if sessiontoken == "" {
				auth.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "missing session token",
				})
				return
			}

			session, err := authService.ValidateSession(sessiontoken)
			if err != nil {
				auth.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid session token",
				})
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				session.UserID,
			)
			ctx = context.WithValue(
				ctx,
				SessionIDKey,
				session.ID,
			)

			r = r.WithContext(ctx)


			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetSessionID(ctx context.Context) (string, bool) {
	SessionID, ok := ctx.Value(SessionIDKey).(string)
	return SessionID, ok 
}
