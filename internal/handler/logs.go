package handler

import "net/http"

func HandleLogs(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error" : "method not allowed",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"log" : "ok",
	})
}