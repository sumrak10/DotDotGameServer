package manager

import (
	"OnlineGame/internal/clients"
	"OnlineGame/internal/database"
	"OnlineGame/internal/utils"
	"errors"
	"time"
)

type GameView struct {
	MatchID       uint                 `json:"id"`
	CreatedAt     time.Time            `json:"created_at"`
	Players       []*database.UserView `json:"players"`
	OwnerPlayerID uint                 `json:"owner_player_id"`
	MinPlayers    uint8                `json:"min_players"`
	MaxPlayers    uint8                `json:"max_players"`
	Status        string               `json:"status"`
}

func (m *Manager) GetAllMatches(clientID uint) []*GameView {
	gamesList := make([]*GameView, 0)
	currentMatchID, clientCurrentlyInMatch := m.IsClientInMatch(clientID)
	if clientCurrentlyInMatch {
		_game, isActiveGame := m.activeGames[currentMatchID]
		if isActiveGame {
			gamesList = append(gamesList, &GameView{
				MatchID:       _game.Match.ID,
				CreatedAt:     _game.Match.CreatedAt,
				Players:       m.GetMatchLobbyUsers(currentMatchID),
				OwnerPlayerID: _game.Match.OwnerPlayerID,
				MinPlayers:    _game.ExampleWorld.MinPlayers,
				MaxPlayers:    _game.ExampleWorld.MaxPlayers,
				Status:        "in_match",
			})
		}
		_game, _isIdleGame := m.idleGames[currentMatchID]
		if _isIdleGame {
			gamesList = append(gamesList, &GameView{
				MatchID:       _game.Match.ID,
				CreatedAt:     _game.Match.CreatedAt,
				Players:       m.GetMatchLobbyUsers(currentMatchID),
				OwnerPlayerID: _game.Match.OwnerPlayerID,
				MinPlayers:    _game.ExampleWorld.MinPlayers,
				MaxPlayers:    _game.ExampleWorld.MaxPlayers,
				Status:        "in_match_lobby",
			})
		}
	}
	for _matchID, _game := range m.idleGames {
		if _matchID == currentMatchID {
			continue
		}
		gamesList = append(gamesList, &GameView{
			MatchID:       _game.Match.ID,
			CreatedAt:     _game.Match.CreatedAt,
			Players:       m.GetMatchLobbyUsers(_matchID),
			OwnerPlayerID: _game.Match.OwnerPlayerID,
			MinPlayers:    _game.ExampleWorld.MinPlayers,
			MaxPlayers:    _game.ExampleWorld.MaxPlayers,
			Status:        "idle_match",
		})
	}
	for _matchID, _game := range m.activeGames {
		if _matchID == currentMatchID {
			continue
		}
		gamesList = append(gamesList, &GameView{
			MatchID:       _game.Match.ID,
			CreatedAt:     _game.Match.CreatedAt,
			Players:       m.GetMatchLobbyUsers(_matchID),
			OwnerPlayerID: _game.Match.OwnerPlayerID,
			MinPlayers:    _game.ExampleWorld.MinPlayers,
			MaxPlayers:    _game.ExampleWorld.MaxPlayers,
			Status:        "active_match",
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
	_currentMatchID, _clientInMatch := m.IsClientInMatch(client.User.ID)
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
	m.addToLobby(client, matchID)

	// Notify other clients
	m.notifyPlayerJoinMatchLobby(matchID, client)

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
	_, _clientInMatch := m.IsClientInMatch(client.User.ID)
	if !_clientInMatch {
		return errors.New("you are now not in any match")
	}

	_game, _idleGameFound := m.idleGames[matchID]
	_, _activeGameFound := m.activeGames[matchID]
	if !_idleGameFound && !_activeGameFound {
		return errors.New("game does not exist")
	}
	if !_idleGameFound && _activeGameFound {
		return errors.New("game is running. wait until it ends or surrender before")
	}

	// Change match owner if player was owner
	if _game.Match.OwnerPlayerID == client.User.ID {
		_game.Match.OwnerPlayerID = utils.ListRandomElement(m.lobby[matchID]).User.ID
	}

	m.removeFromLobby(client.User.ID, matchID)

	if len(m.lobby[matchID]) == 0 {
		// Delete empty lobby game
		m.deleteLobby(matchID)
		m.deleteIdleGame(matchID)
	} else {
		// Notify other clients
		m.notifyPlayerLeaveMatchLobby(matchID, client)
	}
	return nil
}
