package usersAPI

import (
	"OnlineGame/internal/database"
	authAPI "OnlineGame/internal/server/api/auth"
	APIhelpers "OnlineGame/internal/server/api/helpers"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/users/me", GetMe).Methods("GET")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userID, _authorized := authAPI.GetUserIDFromRequest(r)
	if !_authorized {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := database.NewUserRepository(database.GetDB())
	found, err := repo.FindByID(userID)
	if err != nil || found == nil {
		APIhelpers.ErrorJSONResponse(w, "user not found", http.StatusNotFound)
		return
	}
	APIhelpers.JSONResponse(w, found.ToSensitiveView(), http.StatusOK)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, _authorized := authAPI.GetUserIDFromRequest(r)
	if !_authorized {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		APIhelpers.ErrorJSONResponse(w, "invalid body", http.StatusBadRequest)
		return
	}

	repo := database.NewUserRepository(database.GetDB())
	updated, err := repo.Update(userID, updates)
	if err != nil || updated == nil {
		APIhelpers.ErrorJSONResponse(w, "failed to update user", http.StatusInternalServerError)
		return
	}
	APIhelpers.JSONResponse(w, updated, http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, _authorized := authAPI.GetUserIDFromRequest(r)
	if !_authorized {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := database.NewUserRepository(database.GetDB())
	if err := repo.Delete(userID); err != nil {
		APIhelpers.ErrorJSONResponse(w, "failed to delete user", http.StatusInternalServerError)
		return
	}
	APIhelpers.JSONResponse(w, map[string]string{"message": "deleted"}, http.StatusOK)
}
