package manager

import (
	"OnlineGame/clients"
	"OnlineGame/game"
)

func (m *Manager) notifyPlayerJoinMatchLobby(matchID uint, joinedClient *clients.Client) {
	for _, client := range m.lobby[matchID] {
		client.Send(game.NotificationMessage{
			Type:           "match_lobby_join",
			Message:        "Player join lobby",
			MatchID:        matchID,
			PlayerID:       joinedClient.User.ID,
			PlayerUserName: joinedClient.User.UserName,
		})
	}
}

func (m *Manager) notifyPlayerLeaveMatchLobby(matchID uint, joinedClient *clients.Client) {
	for _, client := range m.lobby[matchID] {
		client.Send(game.NotificationMessage{
			Type:           "match_lobby_leave",
			Message:        "Player leave lobby",
			MatchID:        matchID,
			PlayerID:       joinedClient.User.ID,
			PlayerUserName: joinedClient.User.UserName,
		})
	}
}
