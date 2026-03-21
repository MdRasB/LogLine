package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"encoding/json"
)

type HeathResponse struct {
	Status string `json:"status"`

}

type IngestResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	logServe := http.NewServeMux()

	logServe.HandleFunc("GET /health", handleHealth)
	logServe.HandleFunc("POST /ingest", handleIngest)


	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", logServe)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request){
	//msg := HeathResponse{}
	if r.Method != http.MethodGet {
		writeJsonResponse(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	writeJsonResponse(w, http.StatusOK, HeathResponse{
		Status: "ok",
	})
}


func handleIngest(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		writeJsonResponse(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	writeJsonResponse(w,http.StatusOK, IngestResponse{
		Message: "received",
	})

}

func writeJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}