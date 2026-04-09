package APIhelpers

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func CreatedJsonResponse(w http.ResponseWriter, createdID uint) {
	JSONResponse(w, map[string]uint{"id": createdID}, http.StatusCreated)
}

func SuccessJSONResponse(w http.ResponseWriter) {
	JSONResponse(w, map[string]string{"status": "success"}, http.StatusOK)
}

func ErrorJSONResponse(w http.ResponseWriter, message string, status int) {
	JSONResponse(w, map[string]string{"status": "error", "details": message}, status)
}
