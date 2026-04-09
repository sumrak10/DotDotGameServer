package matchesAPI

import (
	"OnlineGame/database"
	"OnlineGame/manager"
	authAPI "OnlineGame/server/api/auth"
	APIhelpers "OnlineGame/server/api/helpers"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RegisterMatchRoutes(r *mux.Router) {
	r.HandleFunc("/matches", CreateMatch).Methods("POST")
	//r.HandleFunc("/matches", GetIdleMatches).Methods("GET")
	r.HandleFunc("/matches/{match_id}/start", StartMatch).Methods("POST")
	//r.HandleFunc("/matches/{match_id}/surrender", SurrenderMatch).Methods("POST")

	r.HandleFunc("/matches/{match_id}/lobby/join", JoinMatchLobby).Methods("POST")
	r.HandleFunc("/matches/{match_id}/lobby/leave", LeaveMatchLobby).Methods("POST")
}

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	userID, _authorized := authAPI.GetUserIDFromRequest(r)
	if !_authorized {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	m := manager.GetManager()
	if m.IsClientInMatch(userID) {
		APIhelpers.ErrorJSONResponse(w, "user currently in match", http.StatusBadRequest)
		return
	}

	// Action
	match := database.Match{}
	database.GetDB().Create(&match)

	err := m.CreateGame(userID, &match)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	APIhelpers.CreatedJsonResponse(w, match.ID)
}

func StartMatch(w http.ResponseWriter, r *http.Request) {
	userID, _authorized := authAPI.GetUserIDFromRequest(r)
	if !_authorized {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	m := manager.GetManager()
	_vars := mux.Vars(r)

	// MatchID Parse
	_matchID, err := strconv.Atoi(_vars["match_id"])
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	matchID := uint(_matchID)

	// Action
	err = m.StartGame(userID, matchID)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	APIhelpers.SuccessJSONResponse(w)
}

func JoinMatchLobby(w http.ResponseWriter, r *http.Request) {
	userID, _authorized := authAPI.GetUserIDFromRequest(r)
	if !_authorized {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	m := manager.GetManager()
	_vars := mux.Vars(r)

	// MatchID Parse
	_matchID, err := strconv.Atoi(_vars["match_id"])
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	matchID := uint(_matchID)

	// Action
	err = m.JoinMatchLobby(userID, matchID)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby := m.GetMatchLobbyUsers(matchID)
	lobbyView := make([]database.UserView, 0, len(lobby))
	for _, user := range lobby {
		if user != nil {
			lobbyView = append(lobbyView, user.ToView())
		}
	}

	APIhelpers.JSONResponse(w, lobbyView, http.StatusOK)
}

func LeaveMatchLobby(w http.ResponseWriter, r *http.Request) {
	userID, ok := authAPI.GetUserIDFromRequest(r)
	if !ok {
		APIhelpers.ErrorJSONResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	m := manager.GetManager()
	_vars := mux.Vars(r)

	// MatchID Parse
	_matchID, err := strconv.Atoi(_vars["match_id"])
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	matchID := uint(_matchID)

	// Action
	err = m.LeaveMatchLobby(userID, matchID)
	if err != nil {
		APIhelpers.ErrorJSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	APIhelpers.SuccessJSONResponse(w)
}
