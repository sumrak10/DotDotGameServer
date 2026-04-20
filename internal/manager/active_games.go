package manager

import "OnlineGame/internal/game"

func (m *Manager) addActiveGame(_game *game.Game) {
	m.activeGamesMu.Lock()
	m.activeGames[_game.Match.ID] = _game
	m.activeGamesMu.Unlock()
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
