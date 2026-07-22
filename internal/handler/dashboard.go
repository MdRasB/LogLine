package handler

import (
	"net/http"

	"github.com/MdRasB/LogLine/internal/dashboard"
	"github.com/MdRasB/LogLine/internal/db"
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
	data := dashboard.DashboardData{}

	err := h.templates.Render(
		w,
		"dashboard.html",
		data,
	)
	if err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
	}
}
func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {}
