package server

import (
	"OnlineGame/internal/manager"
	authAPI "OnlineGame/internal/server/api/auth"
	htmlAPI "OnlineGame/internal/server/api/html"
	matchesAPI "OnlineGame/internal/server/api/matches"
	usersAPI "OnlineGame/internal/server/api/users"
	wsAPI "OnlineGame/internal/server/api/ws"
	staticFiles "OnlineGame/proto"
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

	// Static
	fileServer := http.FileServer(http.FS(staticFiles.StaticFiles))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

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
