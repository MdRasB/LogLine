package handler

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	WriteJSON(w, http.StatusOK, HealthResponse{
		Status: "ok",
	})
}