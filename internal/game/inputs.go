package game

import (
	"OnlineGame/internal/game/world/nodes"
	gamepb "OnlineGame/pkg/pb/go/game"
	"fmt"
)

type inputQueueMessage struct {
	ClientID uint
	Action   *gamepb.ActionInputMessage
}

func (g *Game) HandleInput(clientID uint, action *gamepb.ActionInputMessage) {
	g.inputQueue <- &inputQueueMessage{
		ClientID: clientID,
		Action:   action,
	}
}

func (g *Game) processInputs() {
	for {
		select {
		case input := <-g.inputQueue:
			g.applyInput(input)
		default:
			return
		}
	}
}

func (g *Game) applyInput(queueMessage *inputQueueMessage) {
	player := g.players[queueMessage.ClientID]

	var actionErr error
	switch action := queueMessage.Action.Payload.(type) {
	case *gamepb.ActionInputMessage_SendArmy:
		actionErr = g.world.SendArmy(
			player.Client.User.ID,
			nodes.NodeID(action.SendArmy.HeadingFromId),
			nodes.NodeID(action.SendArmy.HeadingToId),
			uint(action.SendArmy.Value),
		)
	case *gamepb.ActionInputMessage_UpdateNodeType:
		actionErr = g.world.UpdateNodeType(
			player.Client.User.ID,
			nodes.NodeID(action.UpdateNodeType.NodeId),
			action.UpdateNodeType.NewType,
		)
	case *gamepb.ActionInputMessage_SetAlwaysSendArmy:
		actionErr = g.world.SetAlwaysSendArmy(
			player.Client.User.ID,
			nodes.NodeID(action.SetAlwaysSendArmy.FromNodeId),
			nodes.NodeID(action.SetAlwaysSendArmy.ToNodeId),
			action.SetAlwaysSendArmy.Mode,
		)
	default:
		actionErr = fmt.Errorf("unknown action type: %T", action)
	}
	g.sendActionOutputMessage(player.Client, actionErr)
}
