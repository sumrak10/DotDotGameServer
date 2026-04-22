package html

import (
	"OnlineGame/internal/config"
	"OnlineGame/templates"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHtmlRoutes(r *mux.Router) {
	r.HandleFunc("/", Index).Methods("GET")
}

func Index(w http.ResponseWriter, r *http.Request) {
	context := struct {
		BaseURL         string
		ValuesScaleCoef uint
	}{
		BaseURL:         config.Server().GetAddress(),
		ValuesScaleCoef: config.Game().ValuesScaleCoef,
	}
	tmpl, err := template.ParseFS(templates.TemplatesFS, "index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = tmpl.Execute(w, context)
	if err != nil {
		return
	}
}
