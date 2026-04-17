package manager

import (
	"OnlineGame/clients"
)

func (m *Manager) AddClient(client *clients.Client) {
	_inGame, _currentMatchID := m.IsClientInMatch(client.User.ID)
	if _inGame {
		_game, _isActiveGame := m.activeGames[_currentMatchID]
		if _isActiveGame {
			_game.ChangePlayerClient(client)
			_game.NotifyPlayerReconnnected(client.User.ID)
		}
	}
	m.clientsMu.Lock()
	m.clients[client.User.ID] = client
	m.clientsMu.Unlock()
}

func (m *Manager) RemoveClient(clientID uint) {
	m.clientsMu.Lock()
	client, ok := m.clients[clientID]
	client.Conn = nil
	if ok {
		delete(m.clients, clientID)
	}
	m.clientsMu.Unlock()

	_inGame, _currentMatchID := m.IsClientInMatch(clientID)
	if _inGame {
		_activeGame, _isActiveGame := m.activeGames[_currentMatchID]
		if _isActiveGame {
			_activeGame.NotifyPlayerDisconnnected(clientID)
		}
	}
}

func (m *Manager) IsClientInMatch(clientID uint) (bool, uint) {
	matchID, clientInMatch := m.clientGame[clientID]
	if clientInMatch {
		return true, matchID
	}
	return false, 0
}
