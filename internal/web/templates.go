// Package web manages HTML template rendering.
package web

import (
	"html/template"
	"net/http"
	"time"
)

type TemplateManager struct {
	templates *template.Template
}

func NewTemplateManager() (*TemplateManager, error) {
	funcMap := template.FuncMap{
		"add":        add,
		"sub":        sub,
		"formatTime": formatTime,
	}

	tmpl := template.New("").Funcs(funcMap)

	var err error

	tmpl, err = tmpl.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}

	_, err = tmpl.ParseGlob("web/templates/partials/*.html")
	if err != nil {
		return nil, err
	}

	return &TemplateManager{
		templates: tmpl,
	}, nil
}

func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data any) error {
	return tm.templates.ExecuteTemplate(
		w,
		name,
		data,
	)
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func formatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t
	}

	return parsed.Format("02 Jan 2006 15:04:05")
}
