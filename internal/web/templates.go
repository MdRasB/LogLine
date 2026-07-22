// Package web manages HTML template rendering.
package web

import (
	"html/template"
	"net/http"
)

type TemplateManager struct {
	templates *template.Template
}

func NewTemplateManager() (*TemplateManager, error) {
	tmpl, err := template.ParseGlob("web/templates/*.html")
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

func (tm *TemplateManager) Render(
	w http.ResponseWriter,
	name string,
	data any,
) error {
	return tm.templates.ExecuteTemplate(
		w,
		name,
		data,
	)
}
