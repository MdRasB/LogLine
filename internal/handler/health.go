package handler

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	writeJSON(w, http.StatusOK, HealthResponse{
		Status: "ok",
	})
}