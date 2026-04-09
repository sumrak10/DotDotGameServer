package manager

import (
	"OnlineGame/clients"
)

func (m *Manager) AddClient(client *clients.Client) {
	if m.IsClientInMatch(client.User.ID) {
		_game := m.activeGames[client.User.ID]
		_game.NotifyPlayerReconnnected(client.User.ID)
	}
	m.clientsMu.Lock()
	m.clients[client.User.ID] = client
	m.clientsMu.Unlock()
}

func (m *Manager) RemoveClient(clientID uint) {
	if m.IsClientInMatch(clientID) {
		_game := m.activeGames[clientID]
		_game.NotifyPlayerDisconnnected(clientID)
	}
	m.clientsMu.Lock()
	_, ok := m.clients[clientID]
	if ok {
		delete(m.clients, clientID)
	}
	m.clientsMu.Unlock()
}

func (m *Manager) IsClientInMatch(clientID uint) bool {
	_, clientInMatch := m.clientGame[clientID]
	return clientInMatch
}
