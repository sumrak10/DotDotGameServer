package manager

import (
	"OnlineGame/clients"
	"OnlineGame/database"
	"errors"
	"fmt"
	"time"
)

type GameView struct {
	MatchID    uint                 `json:"id"`
	CreatedAt  time.Time            `json:"created_at"`
	Players    []*database.UserView `json:"players"`
	Admin      *database.UserView   `json:"owner"`
	MinPlayers uint8                `json:"min_players"`
	MaxPlayers uint8                `json:"max_players"`
	Status     string               `json:"status"`
}

func (m *Manager) GetAllMatches(clientID uint) []*GameView {
	gamesList := make([]*GameView, 0)
	currentMatchID, clientCurrentlyInMatch := m.clientGame[clientID]
	if clientCurrentlyInMatch {
		_game, isActiveGame := m.activeGames[currentMatchID]
		if isActiveGame {
			gamesList = append(gamesList, &GameView{
				MatchID:    _game.Match.ID,
				CreatedAt:  _game.Match.CreatedAt,
				Players:    m.GetMatchLobbyUsers(currentMatchID),
				Admin:      m.clients[_game.OwnerPlayerID()].User.ToView(),
				MinPlayers: _game.ExampleWorld.MinPlayers,
				MaxPlayers: _game.ExampleWorld.MaxPlayers,
				Status:     "in_match",
			})
		}
		_game, _isIdleGame := m.idleGames[currentMatchID]
		if _isIdleGame {
			gamesList = append(gamesList, &GameView{
				MatchID:    _game.Match.ID,
				CreatedAt:  _game.Match.CreatedAt,
				Players:    m.GetMatchLobbyUsers(currentMatchID),
				Admin:      m.clients[_game.OwnerPlayerID()].User.ToView(),
				MinPlayers: _game.ExampleWorld.MinPlayers,
				MaxPlayers: _game.ExampleWorld.MaxPlayers,
				Status:     "in_match_lobby",
			})
		}
	}
	for _matchID, _game := range m.idleGames {
		if _matchID == currentMatchID {
			continue
		}
		gamesList = append(gamesList, &GameView{
			MatchID:    _game.Match.ID,
			CreatedAt:  _game.Match.CreatedAt,
			Players:    m.GetMatchLobbyUsers(_matchID),
			Admin:      m.clients[_game.OwnerPlayerID()].User.ToView(),
			MinPlayers: _game.ExampleWorld.MinPlayers,
			MaxPlayers: _game.ExampleWorld.MaxPlayers,
			Status:     "idle_match",
		})
	}
	for _matchID, _game := range m.activeGames {
		if _matchID == currentMatchID {
			continue
		}
		gamesList = append(gamesList, &GameView{
			MatchID:    _game.Match.ID,
			CreatedAt:  _game.Match.CreatedAt,
			Players:    m.GetMatchLobbyUsers(_matchID),
			Admin:      m.clients[_game.OwnerPlayerID()].User.ToView(),
			MinPlayers: _game.ExampleWorld.MinPlayers,
			MaxPlayers: _game.ExampleWorld.MaxPlayers,
			Status:     "active_match",
		})
	}
	return gamesList
}

func (m *Manager) GetMatchLobbyUsers(matchID uint) []*database.UserView {
	allClients := m.lobby[matchID]
	users := make([]*database.UserView, 0)
	for _, client := range allClients {
		users = append(users, client.User.ToView())
	}
	fmt.Println(users)
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

	// Add to match lobby
	m.addToMatchLobby(client, matchID)

	// Notify other clients
	m.notifyPlayerJoinMatchLobby(matchID, client)

	return nil
}

func (m *Manager) addToMatchLobby(client *clients.Client, matchID uint) {
	// Add client to lobby
	m.lobbyMu.Lock()
	m.lobby[matchID] = append(m.lobby[matchID], client)
	m.lobbyMu.Unlock()
	// Add client to clientGame
	m.clientGameMu.Lock()
	m.clientGame[client.User.ID] = matchID
	m.clientGameMu.Unlock()
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

	m.removeFromLobby(client, matchID)

	if len(m.lobby[matchID]) == 0 {
		// Delete empty lobby game
		m.deleteIdleGame(matchID)
	} else {
		// Notify other clients
		m.notifyPlayerLeaveMatchLobby(matchID, client)
	}
	return nil
}

func (m *Manager) removeFromLobby(client *clients.Client, matchID uint) {
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
}
