package ws

import (
	"OnlineGame/internal/clients"
	"OnlineGame/internal/database"
	"OnlineGame/internal/manager"
	authAPI "OnlineGame/internal/server/api/auth"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterWSRoutes(r *mux.Router) {
	r.HandleFunc("/ws", HandleWS)
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	// Auth
	_tokenStr := r.URL.Query().Get("token")
	userID, err := authAPI.ParseToken(_tokenStr)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := database.NewUserRepository(database.GetDB())
	user, err := repo.FindByID(userID)
	if err != nil || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	// Manager
	m := manager.GetManager()

	// Upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	// Provide manager
	client := clients.NewClient(user, conn)
	m.AddClient(client)

	go client.WriteLoop()
	client.ReadLoop(
		func(userID uint, data []byte) {
			m.OnMessage(userID, data)
		},
	)

	err = client.Conn.Close()
	if err != nil {
		panic(err)
	}

	m.RemoveClient(userID)
}
