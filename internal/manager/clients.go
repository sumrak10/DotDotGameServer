package manager

import (
	"OnlineGame/internal/clients"
	gamepb "OnlineGame/pkg/pb/go/game"
)

func (m *Manager) AddClient(client *clients.Client) {
	_currentMatchID, _inGame := m.IsClientInMatch(client.User.ID)
	if _inGame {
		_game, _isActiveGame := m.activeGames[_currentMatchID]
		if _isActiveGame {
			_game.ChangePlayerClient(client)
			m.broadcastMessage(_currentMatchID, &gamepb.OutputMessage{
				Payload: &gamepb.OutputMessage_Notification{
					Notification: &gamepb.Notification{
						Type:     "player_reconnected",
						Message:  "Player reconnected",
						MatchId:  uint64(_currentMatchID),
						PlayerId: uint64(client.User.ID),
					},
				},
			})
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

	_currentMatchID, _inGame := m.IsClientInMatch(client.User.ID)
	if _inGame {
		_, _isActiveGame := m.activeGames[_currentMatchID]
		if _isActiveGame {
			m.broadcastMessage(_currentMatchID, &gamepb.OutputMessage{
				Payload: &gamepb.OutputMessage_Notification{
					Notification: &gamepb.Notification{
						Type:     "player_disconnected",
						Message:  "Player disconnected",
						MatchId:  uint64(_currentMatchID),
						PlayerId: uint64(clientID),
					},
				},
			})
		}
	}
}
