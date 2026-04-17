package manager

import (
	"OnlineGame/clients"
	"OnlineGame/game"
	"encoding/json"
)

func (m *Manager) OnMessage(clientID uint, data []byte) {
	client, clientFound := m.clients[clientID]
	if !clientFound {
		panic("client not found")
		return
	}
	var inputMessage game.InputMessage
	err := json.Unmarshal(data, &inputMessage)
	if err != nil {
		panic(err)
	}
	switch inputMessage.Type {
	case "join_lobby":
		m.onJoinMatchMessage(&inputMessage, client)
	case "leave_lobby":
		m.onLeaveMatchMessage(&inputMessage, client)
	default:
		matchID, matchFound := m.clientGame[clientID]
		if !matchFound {
			client.Send("match not found")
			return
		}
		m.activeGames[matchID].HandleInput(clientID, inputMessage)
	}
}

type LobbyMessagePayload struct {
	MatchID uint `json:"match_id"`
}
type LobbyMessageResponse struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

func (m *Manager) onJoinMatchMessage(inputMessage *game.InputMessage, client *clients.Client) {
	// Input
	var joinMatchPayload LobbyMessagePayload
	err := json.Unmarshal(inputMessage.Payload, &joinMatchPayload)
	if err != nil {
		panic(err)
	}
	// Action
	actionErr := m.joinMatchLobby(client, joinMatchPayload.MatchID)
	// Output
	game.SendActionMessageResponse(client, inputMessage.Type, actionErr)
}

func (m *Manager) onLeaveMatchMessage(inputMessage *game.InputMessage, client *clients.Client) {
	// Input
	var leaveMatchPayload LobbyMessagePayload
	err := json.Unmarshal(inputMessage.Payload, &leaveMatchPayload)
	if err != nil {
		panic(err)
	}
	// Action
	actionErr := m.leaveMatchLobby(client, leaveMatchPayload.MatchID)
	// Output
	game.SendActionMessageResponse(client, inputMessage.Type, actionErr)
}
