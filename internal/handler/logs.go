package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MdRasB/LogLine/internal/db"
	"github.com/MdRasB/LogLine/internal/model"
)


type LogHandler struct {
	store *db.DBStore
}

func NewLogHandler(store *db.DBStore) *LogHandler {
	return &LogHandler{store: store}
}

func (h *LogHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"})
		return
	}

	pf, err := parseFilter(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest,
            map[string]string{"error": err.Error()})
        return
	}

	result, err := h.store.GetLogs(pf, r.Context())
	if err != nil {
        WriteJSON(w, http.StatusInternalServerError,
            map[string]string{"error": "failed to retrieve logs"})
        return
    }

	WriteJSON(w, http.StatusOK, result)

}

var validLevels = map[string]bool{
    "error": true, "warn": true,
    "info":  true, "debug": true, "fatal": true,
}


func parseFilter(r *http.Request) (model.LogFilter, error) {
    q := r.URL.Query()

    // page — default 1, must be >= 1
    page, err := strconv.Atoi(q.Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    // limit — default 20, hard cap at 100 (prevents OOM attacks)
    limit, err := strconv.Atoi(q.Get("limit"))
    if err != nil || limit < 1 || limit > 100 {
        limit = 20
    }

    // level — must be in whitelist, empty string means "no filter"
    level := q.Get("level")
    if level != "" && !validLevels[level] {
        return model.LogFilter{}, fmt.Errorf(
            "invalid level %q — must be: error, warn, info, debug, fatal", level)
    }

    // from / to — optional RFC3339 timestamps
    from, err := parseTimeParam(q.Get("from"))
    if err != nil {
        return model.LogFilter{}, fmt.Errorf("invalid 'from': %w", err)
    }

    to, err := parseTimeParam(q.Get("to"))
    if err != nil {
        return model.LogFilter{}, fmt.Errorf("invalid 'to': %w", err)
    }

    return model.LogFilter{
        Service: q.Get("service"),
        Level:   level,
        Search:  q.Get("search"),
        From:    from,
        To:      to,
        Page:    page,
        Limit:   limit,
    }, nil
}

func parseTimeParam(val string) (*time.Time, error) {
    if val == "" {
        return nil, nil
    }
    t, err := time.Parse(time.RFC3339, val)
    if err != nil {
        return nil, fmt.Errorf("use RFC3339 format e.g. 2024-01-15T10:00:00Z")
    }
    utc := t.UTC()
    return &utc, nil
}