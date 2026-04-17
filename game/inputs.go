package game

import (
	"OnlineGame/game/world/nodes"
	"encoding/json"
	"fmt"
)

type PlayerInput struct {
	ClientID     uint
	InputMessage InputMessage
}

func (g *Game) HandleInput(clientID uint, inputMessage InputMessage) {
	g.inputQueue <- PlayerInput{
		ClientID:     clientID,
		InputMessage: inputMessage,
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

func (g *Game) applyInput(playerInput PlayerInput) {
	player := g.players[playerInput.ClientID]
	switch playerInput.InputMessage.Type {
	case "send_army_action":
		var sendArmyMessage SendArmyMessagePayload
		err := json.Unmarshal(playerInput.InputMessage.Payload, &sendArmyMessage)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received action", sendArmyMessage)
		g.sendArmyAction(
			player,
			playerInput.InputMessage.Type,
			sendArmyMessage.HeadingFromID,
			sendArmyMessage.HeadingToID,
			sendArmyMessage.Value,
		)
	case "update_node_type_action":
		var updateNodeTypeMessage UpdateNodeTypeMessagePayload
		err := json.Unmarshal(playerInput.InputMessage.Payload, &updateNodeTypeMessage)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received action", updateNodeTypeMessage)
		g.updateNodeTypeAction(
			player,
			playerInput.InputMessage.Type,
			updateNodeTypeMessage.NodeID,
			updateNodeTypeMessage.NewType,
		)
	case "set_always_send_army_action":
		var setAlwaysSendArmyMessagePayload SetAlwaysSendArmyMessagePayload
		err := json.Unmarshal(playerInput.InputMessage.Payload, &setAlwaysSendArmyMessagePayload)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received action", setAlwaysSendArmyMessagePayload)
		g.setAlwaysSendArmyAction(
			player,
			playerInput.InputMessage.Type,
			setAlwaysSendArmyMessagePayload.FromNodeID,
			setAlwaysSendArmyMessagePayload.ToNodeID,
			setAlwaysSendArmyMessagePayload.Mode,
		)
	}
}

func (g *Game) sendArmyAction(player *Player, inputMessageType string, headingFromID nodes.NodeID, headingToID nodes.NodeID, Value uint) {
	actionErr := g.world.SendArmy(player.Client.User.ID, headingFromID, headingToID, Value)
	// Output
	SendActionMessageResponse(player.Client, inputMessageType, actionErr)
}

func (g *Game) updateNodeTypeAction(player *Player, inputMessageType string, nodeID nodes.NodeID, NewType nodes.NodeType) {
	actionErr := g.world.UpdateNodeType(player.Client.User.ID, nodeID, NewType)
	// Output
	SendActionMessageResponse(player.Client, inputMessageType, actionErr)
}

func (g *Game) setAlwaysSendArmyAction(player *Player, inputMessageType string, fromNodeID nodes.NodeID, toNodeID nodes.NodeID, mode bool) {
	actionErr := g.world.SetAlwaysSendArmy(player.Client.User.ID, fromNodeID, toNodeID, mode)
	// Output
	SendActionMessageResponse(player.Client, inputMessageType, actionErr)
}
