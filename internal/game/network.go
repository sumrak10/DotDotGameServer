package game

import (
	"OnlineGame/internal/clients"
	gamepb "OnlineGame/pkg/pb/go/game"
)

func (g *Game) broadcastMessage(msg *gamepb.OutputMessage) {
	for _, player := range g.players {
		if player.Client.Conn != nil {
			player.Client.Send(msg)
		}
	}
}

func (g *Game) sendOutputMessage(client *clients.Client, payloadResponse *gamepb.OutputMessage_ActionOutputMessage) {
	client.Send(&gamepb.OutputMessage{Payload: payloadResponse})
}

func (g *Game) sendActionOutputMessage(client *clients.Client, actionErr error) {
	if actionErr == nil {
		g.sendOutputMessage(client, &gamepb.OutputMessage_ActionOutputMessage{
			ActionOutputMessage: &gamepb.ActionOutputMessage{
				Status: "success",
			},
		})
	} else {
		g.sendOutputMessage(client, &gamepb.OutputMessage_ActionOutputMessage{
			ActionOutputMessage: &gamepb.ActionOutputMessage{
				Status:  "error",
				Details: actionErr.Error(),
			},
		})
	}

}
