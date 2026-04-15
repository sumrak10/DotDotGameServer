package manager

import (
	"OnlineGame/clients"
	"OnlineGame/game"
	"sync"
)

var (
	instance *Manager
	once     sync.Once
)

type Manager struct {
	idleGamesMu sync.RWMutex
	idleGames   map[uint]*game.Game // MatchID - Game

	activeGamesMu sync.RWMutex
	activeGames   map[uint]*game.Game // MatchID - Game

	clientsMu sync.RWMutex
	clients   map[uint]*clients.Client // clientID - Client

	clientGameMu sync.RWMutex
	clientGame   map[uint]uint // clientID - MatchID

	lobbyMu sync.RWMutex
	lobby   map[uint][]*clients.Client // MatchID - []Client
}

func NewManager() *Manager {
	return &Manager{
		idleGames:   make(map[uint]*game.Game),
		activeGames: make(map[uint]*game.Game),
		clients:     make(map[uint]*clients.Client),
		clientGame:  make(map[uint]uint),
		lobby:       make(map[uint][]*clients.Client),
	}
}

func GetManager() *Manager {
	once.Do(func() {
		instance = NewManager()
	})
	return instance
}
