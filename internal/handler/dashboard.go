package handler

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/MdRasB/LogLine/internal/dashboard"
	"github.com/MdRasB/LogLine/internal/db"
	"github.com/MdRasB/LogLine/internal/model"
	"github.com/MdRasB/LogLine/internal/web"
)

type DashboardHandler struct {
	logStore  *db.DBStore
	templates *web.TemplateManager
}

func NewDashboardHandler(store *db.DBStore, templates *web.TemplateManager) *DashboardHandler {
	return &DashboardHandler{
		logStore:  store,
		templates: templates,
	}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	filter := h.buildLogFilter(r)

	result, err := h.logStore.GetLogs(
		filter,
		context.Background(),
	)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load dashboard",
		})
		return
	}

	stats, err := h.logStore.GetDashboardStats(
		context.Background(),
	)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load dashboard statistics",
		})
		return
	}

	data := h.buildDashboardData(
		result,
		stats,
		filter,
	)

	err = h.templates.Render(
		w,
		"dashboard.html",
		data,
	)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}
}

func (h *DashboardHandler) buildLogFilter(r *http.Request) model.LogFilter {
	query := r.URL.Query()

	page := 1
	limit := 20

	if p := query.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if l := query.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	return model.LogFilter{
		Service: query.Get("service"),
		Level:   query.Get("level"),
		Search:  query.Get("search"),
		Page:    page,
		Limit:   limit,
	}
}

func (h *DashboardHandler) buildDashboardData(result *model.PaginatedLogs, stats *dashboard.DashboardStats, filter model.LogFilter) dashboard.DashboardData {
	totalPages := int(
		math.Ceil(
			float64(result.Total) /
				float64(result.Limit),
		),
	)

	if totalPages == 0 {
		totalPages = 1
	}

	return dashboard.DashboardData{
		Logs: result.Logs,

		Stats: *stats,

		Filters: dashboard.DashboardFilters{
			Service: filter.Service,
			Level:   filter.Level,
			Search:  filter.Search,
		},

		Pagination: dashboard.PaginationData{
			CurrentPage: result.Page,
			TotalPages:  totalPages,
			HasNext:     result.HasMore,
			HasPrev:     result.Page > 1,
		},
	}
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
}
