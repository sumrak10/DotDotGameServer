package html

import (
	"OnlineGame/internal/config"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHtmlRoutes(r *mux.Router) {
	r.HandleFunc("/", Index).Methods("GET")
}

func Index(w http.ResponseWriter, r *http.Request) {
	context := struct {
		BaseURL string
	}{
		BaseURL: config.Server().GetAddress(),
	}
	var tmplIndex, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	if tmplIndex == nil {
		fmt.Println("tmplIndex is nil")
		return
	}
	err = tmplIndex.Execute(w, context)
	if err != nil {
		return
	}
}
