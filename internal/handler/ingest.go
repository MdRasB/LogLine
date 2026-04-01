package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MdRasB/LogLine/internal/model"
)

type IngestResponse struct {
	Message string `json:"message"`
}

func HandleIngest(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	var log model.Logs

	err := jsonDecode(r.Body, &log)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err = validate(log); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response := map[string]string{
		"message": "log accepted",
	}

	w.Header().Set("Content_type:", "application/JSON")
	json.NewEncoder(w).Encode(response)

	fmt.Println("Received request:", r.Method)
}

func jsonDecode(body io.Reader, log *model.Logs) error {

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&log); err != nil {
		//http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return errors.New("invalid json")
	}

	return nil

}

func validate(log model.Logs) error {

	if log.Level == "" {
		return errors.New("level is required")
	}

	validLevels := map[string]bool{
		"error": true,
		"debug": true,
		"info":  true,
		"fatal": true,
		"warn":  true,
	}

	if !validLevels[log.Level] {
		return errors.New("invalid level")
	}

	if log.Message == "" {
		return errors.New("message is required")
	}

	if log.Service == "" {
		return errors.New("service is required")
	}

	if log.Timestamp == "" {
		return errors.New("timestamp is required")
	}

	_, err := time.Parse(time.RFC3339, log.Timestamp)
	if err != nil {
		return errors.New("invalid timestamp")
	}

	return nil

}
