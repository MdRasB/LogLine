package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

func GenerateRequestID() (string, error) {
	key := make([]byte, 16)

	if _, err := rand.Read(key); err != nil {

		return "", fmt.Errorf("generate request id: %w", err)
	}

	requestIDkey := "req_" + hex.EncodeToString(key)

	return requestIDkey, nil
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID, err := GenerateRequestID()

		if err != nil {
			http.Error(w,
				"internal server error",
				http.StatusInternalServerError,
			)
		}

		ctx := context.WithValue(
			r.Context(),
			RequestIDKey,
			requestID,
		)

		w.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) (string, bool) {
	RequestID, ok := ctx.Value(RequestIDKey).(string)
	return RequestID, ok
}
