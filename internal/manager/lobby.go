package manager

import "OnlineGame/internal/clients"

func (m *Manager) createEmptyLobby(matchID uint) {
	m.lobbyMu.Lock()
	m.lobby[matchID] = make([]*clients.Client, 0)
	m.lobbyMu.Unlock()
}

func (m *Manager) addToLobby(client *clients.Client, matchID uint) {
	// Add client to lobby
	m.lobbyMu.Lock()
	m.lobby[matchID] = append(m.lobby[matchID], client)
	m.lobbyMu.Unlock()
	// Add client to lobbyClientGame
	m.lobbyClientGameMu.Lock()
	m.lobbyClientGame[client.User.ID] = matchID
	m.lobbyClientGameMu.Unlock()
}

func (m *Manager) removeFromLobby(clientID uint, matchID uint) {
	// Remove client from lobby
	m.lobbyMu.Lock()
	for i, lobbyClient := range m.lobby[matchID] {
		if clientID == lobbyClient.User.ID {
			m.lobby[matchID][i] = nil
			m.lobby[matchID] = append(m.lobby[matchID][:i], m.lobby[matchID][i+1:]...)
		}
	}
	m.lobbyMu.Unlock()
	// Remove client from lobbyClientGame
	m.lobbyClientGameMu.Lock()
	_, ok := m.lobbyClientGame[clientID]
	if ok {
		delete(m.lobbyClientGame, clientID)
	}
	m.lobbyClientGameMu.Unlock()
}

func (m *Manager) deleteLobby(matchID uint) {
	m.lobbyClientGameMu.Lock()
	for _, client := range m.lobby[matchID] {
		delete(m.lobbyClientGame, client.User.ID)
	}
	m.lobbyClientGameMu.Unlock()

	m.lobbyMu.Lock()
	delete(m.lobby, matchID)
	m.lobbyMu.Unlock()
}

func (m *Manager) IsClientInMatch(clientID uint) (uint, bool) {
	matchID, clientInMatch := m.lobbyClientGame[clientID]
	if clientInMatch {
		return matchID, true
	}
	return 0, false
}
