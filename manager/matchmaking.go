package manager

import (
	"OnlineGame/clients"
	"OnlineGame/database"
	"errors"
)

func (m *Manager) GetMatchLobbyUsers(matchID uint) []*database.User {
	allClients := m.lobby[matchID]
	users := make([]*database.User, len(allClients))
	for _, client := range allClients {
		users = append(users, client.User)
	}
	return users
}

func (m *Manager) JoinMatchLobby(clientID uint, matchID uint) error {
	client, clientFound := m.clients[clientID]
	if !clientFound {
		return errors.New("for this action, you should create a websocket connection")
	}
	err := m.joinMatchLobby(client, matchID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) joinMatchLobby(client *clients.Client, matchID uint) error {
	_currentMatchID, _clientInMatch := m.clientGame[client.User.ID]
	if _clientInMatch && _currentMatchID != matchID {
		return errors.New("you are now in other match")
	}
	if _clientInMatch && _currentMatchID == matchID {
		return errors.New("you are now in current match")
	}

	_, _idleGameFound := m.idleGames[matchID]
	_, _activeGameFound := m.activeGames[matchID]
	if !_idleGameFound && !_activeGameFound {
		return errors.New("game does not exist")
	}
	if !_idleGameFound && _activeGameFound {
		return errors.New("game is running")
	}

	// Add client to lobby
	m.lobbyMu.Lock()
	m.lobby[matchID] = append(m.lobby[matchID], client)
	m.lobbyMu.Unlock()
	// Add client to clientGame
	m.clientGameMu.Lock()
	m.clientGame[client.User.ID] = matchID
	m.clientGameMu.Unlock()
	return nil
}

func (m *Manager) LeaveMatchLobby(clientID uint, matchID uint) error {
	client, clientFound := m.clients[clientID]
	if !clientFound {
		return errors.New("for this action, you should create a websocket connection")
	}
	err := m.leaveMatchLobby(client, matchID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) leaveMatchLobby(client *clients.Client, matchID uint) error {
	_, _clientInMatch := m.clientGame[client.User.ID]
	if !_clientInMatch {
		return errors.New("you are now not in any match")
	}

	_, _idleGameFound := m.idleGames[matchID]
	_, _activeGameFound := m.activeGames[matchID]
	if !_idleGameFound && !_activeGameFound {
		return errors.New("game does not exist")
	}
	if !_idleGameFound && _activeGameFound {
		return errors.New("game is running. wait until it ends or surrender before")
	}

	// Remove client from lobby
	m.lobbyMu.Lock()
	for i, lobbyClient := range m.lobby[matchID] {
		if client.User.ID == lobbyClient.User.ID {
			m.lobby[matchID][i] = nil
			m.lobby[matchID] = append(m.lobby[matchID][:i], m.lobby[matchID][i+1:]...)
		}
	}
	m.lobbyMu.Unlock()

	// Remove client from clientGame
	m.clientGameMu.Lock()
	_, ok := m.clientGame[client.User.ID]
	if ok {
		delete(m.clientGame, client.User.ID)
	}
	m.clientGameMu.Unlock()
	return nil
}
