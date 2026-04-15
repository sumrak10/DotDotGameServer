package manager

import (
	"OnlineGame/clients"
	"OnlineGame/database"
	"OnlineGame/game"
	"errors"
	"fmt"
)

func (m *Manager) CreateGame(clientID uint, match *database.Match) error {
	_game, err := game.NewGame(match, "default", clientID)
	if err != nil {
		return err
	}
	_game.OnShutdown = func() {
		_deletedIdleGame := m.deleteIdleGame(match.ID)
		if _deletedIdleGame != nil {
		}
		_deletedActiveGame := m.deleteActiveGame(match.ID)
		if _deletedActiveGame != nil {
		}
	}

	m.idleGamesMu.Lock()
	m.idleGames[_game.Match.ID] = _game
	m.idleGamesMu.Unlock()
	m.lobbyMu.Lock()
	m.lobby[match.ID] = make([]*clients.Client, 0)
	m.lobbyMu.Unlock()

	return nil
}

func (m *Manager) DeleteGame(matchID uint) *game.Game {
	m.lobbyMu.Lock()
	delete(m.lobby, matchID)
	m.lobbyMu.Unlock()
	return m.deleteIdleGame(matchID)
}

func (m *Manager) StartGame(clientID uint, matchID uint) error {
	_game, _foundGame := m.idleGames[matchID]
	if !_foundGame {
		return errors.New("game not in idle status or game does not exist")
	}

	fmt.Println("m.lobby[matchID]")
	fmt.Println(m.lobby[matchID])
	err := _game.Start(clientID, m.lobby[matchID])
	if err != nil {
		return err
	}

	m.deleteIdleGame(matchID)

	m.activeGamesMu.Lock()
	m.activeGames[_game.Match.ID] = _game
	m.activeGamesMu.Unlock()
	return nil
}

func (m *Manager) StopGame(matchID uint) *game.Game {
	_game := m.deleteActiveGame(matchID)
	if _game != nil {
		_game.Stop()
	}
	return _game
}

func (m *Manager) deleteIdleGame(matchID uint) *game.Game {
	m.idleGamesMu.Lock()
	_game, ok := m.idleGames[matchID]
	if ok {
		delete(m.idleGames, matchID)
	}
	m.idleGamesMu.Unlock()
	if ok {
		return _game
	} else {
		return nil
	}
}

func (m *Manager) deleteActiveGame(matchID uint) *game.Game {
	m.activeGamesMu.Lock()
	_game, ok := m.activeGames[matchID]
	if ok {
		delete(m.activeGames, matchID)
	}
	m.activeGamesMu.Unlock()
	if ok {
		return _game
	} else {
		return nil
	}
}
