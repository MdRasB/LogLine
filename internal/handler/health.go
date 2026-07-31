package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Uptime    string `json:"uptime"`
	StartedAt string `json:"started_at"`
	Time      string `json:"time"`
	Version   string `json:"version"`
}

type HealthHandler struct {
	db        *pgxpool.Pool
	startedAt time.Time
	version   string
}

func NewHealthHandler(db *pgxpool.Pool, startedAt time.Time, version string) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startedAt: startedAt,
		version:   version,
	}
}

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	databaseStatus := "up"
	fullStatus := "ok"
	httpStatusCode := http.StatusOK

	// FIX 1: Safely check if database pool is nil BEFORE trying to call .Ping() on it
	if h.db == nil {
		databaseStatus = "down"
		fullStatus = "degraded"
		httpStatusCode = http.StatusServiceUnavailable
	} else {
		// Pool is safe to use, proceed to ping
		err := h.db.Ping(r.Context())
		if err != nil {
			databaseStatus = "down"
			fullStatus = "degraded"
			httpStatusCode = http.StatusServiceUnavailable
		}
	}

	uptimeDuration := time.Since(h.startedAt)
	uptimeStr := uptimeDuration.Truncate(time.Second).String()

	report := HealthResponse{
		Status:    fullStatus,
		Database:  databaseStatus,
		Uptime:    uptimeStr,
		StartedAt: h.startedAt.UTC().Format(time.RFC3339),
		Time:      time.Now().UTC().Format(time.RFC3339),
		Version:   h.version,
	}

	// FIX 2: Pass your actual populated "report" object into WriteJSON
	WriteJSON(w, httpStatusCode, report)
}

