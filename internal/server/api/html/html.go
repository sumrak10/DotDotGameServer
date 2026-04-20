package html

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHtmlRoutes(r *mux.Router) {
	r.HandleFunc("/", Index).Methods("GET")
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		Name: "Пользователь",
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	err := tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

func JsPb(w http.ResponseWriter, r *http.Request) {
}
