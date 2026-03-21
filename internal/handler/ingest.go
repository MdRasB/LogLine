package handler

import (
	"fmt"
	"net/http"
)

type IngestResponse struct {
	Message string `json:"message"`
}

func HandleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	writeJSON(w, http.StatusCreated, IngestResponse{
		Message: "received",
	})

	fmt.Println("Received request:", r.Method)
}