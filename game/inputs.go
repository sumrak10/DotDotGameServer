package game

import (
	"OnlineGame/clients"
	"OnlineGame/game/world/nodes"
	"encoding/json"
)

func (g *Game) ApplyInput(client *clients.Client, inputMessage InputMessage) {
	switch inputMessage.Type {
	case "send_army_action":
		var sendArmyMessage SendArmyMessagePayload
		err := json.Unmarshal(inputMessage.Payload, &sendArmyMessage)
		if err != nil {
			panic(err)
		}
		g.sendArmyAction(
			client,
			inputMessage,
			sendArmyMessage.HeadingFrom,
			sendArmyMessage.HeadingTo,
			sendArmyMessage.Value,
		)
	case "update_node_type_action":
		var updateNodeTypeMessage UpdateNodeTypeMessagePayload
		err := json.Unmarshal(inputMessage.Payload, &updateNodeTypeMessage)
		if err != nil {
			panic(err)
		}
		g.updateNodeTypeAction(
			client,
			inputMessage,
			updateNodeTypeMessage.NodeID,
			updateNodeTypeMessage.NewType,
		)
	case "set_always_send_army_action":
		var setAlwaysSendArmyMessagePayload SetAlwaysSendArmyMessagePayload
		err := json.Unmarshal(inputMessage.Payload, &setAlwaysSendArmyMessagePayload)
		if err != nil {
			panic(err)
		}
		g.setAlwaysSendArmyAction(
			client,
			inputMessage,
			setAlwaysSendArmyMessagePayload.FromNodeID,
			setAlwaysSendArmyMessagePayload.ToNodeID,
			setAlwaysSendArmyMessagePayload.Mode,
		)
	}
}

func (g *Game) sendArmyAction(client *clients.Client, inputMessage InputMessage, headingFrom nodes.NodeID, headingTo nodes.NodeID, Value uint) {
	actionErr := g.world.SendArmy(client.User.ID, headingFrom, headingTo, Value)
	// Output
	SendActionMessageResponse(client, inputMessage.Type, actionErr)
}

func (g *Game) updateNodeTypeAction(client *clients.Client, inputMessage InputMessage, nodeID nodes.NodeID, NewType nodes.NodeType) {
	actionErr := g.world.UpdateNodeType(client.User.ID, nodeID, NewType)
	// Output
	SendActionMessageResponse(client, inputMessage.Type, actionErr)
}

func (g *Game) setAlwaysSendArmyAction(client *clients.Client, inputMessage InputMessage, fromNodeID nodes.NodeID, toNodeID nodes.NodeID, mode bool) {
	actionErr := g.world.SetAlwaysSendArmy(client.User.ID, fromNodeID, toNodeID, mode)
	// Output
	SendActionMessageResponse(client, inputMessage.Type, actionErr)
}
