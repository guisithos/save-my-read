package handlers

import (
	"html/template"
	"net/http"
)

type ViewHandler struct {
	templates *template.Template
}

func NewViewHandler() (*ViewHandler, error) {
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}
	return &ViewHandler{templates: tmpl}, nil
}

func (h *ViewHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "index.html", nil)
}
