package manager

import "OnlineGame/internal/game"

func (m *Manager) addIdleGame(_game *game.Game) {
	m.idleGamesMu.Lock()
	m.idleGames[_game.Match.ID] = _game
	m.idleGamesMu.Unlock()
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
