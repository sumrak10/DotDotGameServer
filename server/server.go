package server

import (
	"OnlineGame/manager"
	authAPI "OnlineGame/server/api/auth"
	htmlAPI "OnlineGame/server/api/html"
	matchesAPI "OnlineGame/server/api/matches"
	usersAPI "OnlineGame/server/api/users"
	wsAPI "OnlineGame/server/api/ws"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	manager    *manager.Manager
	websockets map[*websocket.Conn]uint
}

func NewServer() *Server {
	return &Server{manager: manager.GetManager()}
}

func (s *Server) Start(addr string) error {
	r := mux.NewRouter()

	// Unprotected routes
	htmlAPI.RegisterHtmlRoutes(r)
	authAPI.RegisterAuthRoutes(r)
	wsAPI.RegisterWSRoutes(r)

	// Protected routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(authAPI.Middleware)
	usersAPI.RegisterUserRoutes(protected)
	matchesAPI.RegisterMatchRoutes(protected)

	fmt.Println("Starting server at addr", addr)
	return http.ListenAndServe(addr, r)
}
