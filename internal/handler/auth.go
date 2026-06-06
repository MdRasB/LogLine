package handler

import (
	"encoding/json"
	"errors"
	"strings"

	//"fmt"
	"io"
	"net/http"

	"github.com/MdRasB/LogLine/internal/auth"
	//"github.com/MdRasB/LogLine/internal/db"
	//"github.com/MdRasB/LogLine/internal/model"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	var req auth.RegisterRequest

	if err := jsonDecodeRequest(r.Body, &req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":"registration failed",
		})
		return
	}

	response := map[string]string{
		"message": "registration successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	//fmt.Println("Received request:", r.Method)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1 << 20)

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error":"method not allowed",
		})
		return 
	}

	var login auth.LoginRequest

	if err := jsonDecodeLogin(r.Body, &login); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	sessToken, err := h.authService.Login(login) 
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":"login failed",
		}) 
	}

	response := map[string]string{
		"session_token" : sessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request){
	
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error":"method not allowed",
		})
		return 
	}
	
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing authorization header",
		})
		return
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid authorization format",
		})
		return
	}

	sessionToken := strings.TrimPrefix(authHeader, prefix)

	if sessionToken == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing session token",
		})
		return
	}

	if err := h.authService.Logout(sessionToken); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "logout failed",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func jsonDecodeRequest(body io.Reader, req *auth.RegisterRequest) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil {
		return errors.New("invalid json")
	}

	return nil
}

func jsonDecodeLogin(body io.Reader, req *auth.LoginRequest) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil {
		return errors.New("invalid json")
	}

	return nil 
}
