package manager

import (
	"OnlineGame/internal/clients"
	"OnlineGame/internal/game"
	"sync"
)

var (
	instance *Manager
	once     sync.Once
)

type Manager struct {
	idleGamesMu sync.RWMutex
	idleGames   map[uint]*game.Game // MatchID - Game | Idle games storage

	activeGamesMu sync.RWMutex
	activeGames   map[uint]*game.Game // MatchID - Game | Active games storage

	clientsMu sync.RWMutex
	clients   map[uint]*clients.Client // clientID - Client | Any client who have ws connection

	lobbyMu sync.RWMutex
	lobby   map[uint][]*clients.Client // MatchID - []Client | Match lobby

	lobbyClientGameMu sync.RWMutex
	lobbyClientGame   map[uint]uint // clientID - MatchID | Cache for lobby
}

func NewManager() *Manager {
	return &Manager{
		idleGames:       make(map[uint]*game.Game),
		activeGames:     make(map[uint]*game.Game),
		clients:         make(map[uint]*clients.Client),
		lobbyClientGame: make(map[uint]uint),
		lobby:           make(map[uint][]*clients.Client),
	}
}

func GetManager() *Manager {
	once.Do(func() {
		instance = NewManager()
	})
	return instance
}
