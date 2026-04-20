package manager

import (
	"OnlineGame/internal/clients"
	gamepb "OnlineGame/pkg/pb/go/game"

	"google.golang.org/protobuf/proto"
)

func (m *Manager) OnMessage(clientID uint, data []byte) {
	_, clientFound := m.clients[clientID]
	if !clientFound {
		panic("client not found")
		return
	}

	var inputMessage gamepb.InputMessage
	if err := proto.Unmarshal(data, &inputMessage); err != nil {
		panic("can't unmarshal data")
	}

	switch inputMessage.Payload.(type) {
	case *gamepb.InputMessage_Action:
		matchID, matchFound := m.IsClientInMatch(clientID)
		if !matchFound {
			panic("match not found")
		}
		m.activeGames[matchID].HandleInput(clientID, inputMessage.GetAction())
	default:
		panic("invalid input message")
	}
}

func (m *Manager) broadcastMessage(matchID uint, msg proto.Message) {
	for _, client := range m.lobby[matchID] {
		if client.Conn != nil {
			client.Send(msg)
		}
	}
}

func (m *Manager) notifyPlayerJoinMatchLobby(matchID uint, joinedClient *clients.Client) {
	m.broadcastMessage(matchID, &gamepb.OutputMessage{
		Payload: &gamepb.OutputMessage_Notification{
			Notification: &gamepb.Notification{
				Type:           "match_lobby_join",
				Message:        "Player join lobby",
				MatchId:        uint64(matchID),
				PlayerId:       uint64(joinedClient.User.ID),
				PlayerUsername: joinedClient.User.Username,
			},
		},
	})
}

func (m *Manager) notifyPlayerLeaveMatchLobby(matchID uint, joinedClient *clients.Client) {
	m.broadcastMessage(matchID, &gamepb.OutputMessage{
		Payload: &gamepb.OutputMessage_Notification{
			Notification: &gamepb.Notification{
				Type:           "match_lobby_leave",
				Message:        "Player leave lobby",
				MatchId:        uint64(matchID),
				PlayerId:       uint64(joinedClient.User.ID),
				PlayerUsername: joinedClient.User.Username,
			},
		},
	})
}
