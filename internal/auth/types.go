// Package auth handles the authentication of the servers.
package auth

import (
	//"net/http"
	"encoding/json"
	"net/http"
	"strings"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ExtractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return ""
	}

	sessionToken := strings.TrimPrefix(authHeader, prefix)
	return sessionToken
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
