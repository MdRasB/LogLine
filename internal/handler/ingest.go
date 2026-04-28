package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	//"time"

	"github.com/MdRasB/LogLine/internal/db"
	"github.com/MdRasB/LogLine/internal/model"
)

type IngestHandler struct{
	store *db.DBStore
} 

func NewIngestHandler(store *db.DBStore) *IngestHandler {
	return &IngestHandler{
		store : store,
	}
}

//func HandleIngest(w http.ResponseWriter, r *http.Request) {
func (h *IngestHandler) Handle(w http.ResponseWriter, r *http.Request) { 

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	var log model.Logs

	if err := jsonDecode(r.Body, &log); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := Validate(log); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := h.store.Insert(log); err != nil {
		http.Error(w, "failed to store log", http.StatusInternalServerError)
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
