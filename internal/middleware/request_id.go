package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/MdRasB/LogLine/internal/contextutil"
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
			WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
			return
		}
		ctx := contextutil.WithRequestID(
			r.Context(),
			requestID,
		)

		w.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
