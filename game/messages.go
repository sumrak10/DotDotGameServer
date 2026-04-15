package game

import (
	"OnlineGame/clients"
	"OnlineGame/game/world/nodes"
	"encoding/json"
)

type InputMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendArmyMessagePayload struct {
	Value         uint         `json:"value"`
	HeadingFromID nodes.NodeID `json:"heading_from_id"`
	HeadingToID   nodes.NodeID `json:"heading_to_id"`
}

type UpdateNodeTypeMessagePayload struct {
	NodeID  nodes.NodeID   `json:"node_id"`
	NewType nodes.NodeType `json:"new_type"`
}

type SetAlwaysSendArmyMessagePayload struct {
	FromNodeID nodes.NodeID `json:"from_node_id"`
	ToNodeID   nodes.NodeID `json:"to_node_id"`
	Mode       bool         `json:"mode"`
}

type OutputMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type TickMessagePayload struct {
	MatchID    uint                    `json:"match_id"`
	CurrentTPS float32                 `json:"current_tps"`
	World      WorldTickMessagePayload `json:"world"`
}

type WorldTickMessagePayload struct {
	Nodes     map[nodes.NodeID]*nodes.Node `json:"nodes"`
	NodeEdges []*nodes.NodeEdge            `json:"node_edges"`
}

type ActionMessageResponse struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

func BuildActionMessageResponse(actionErr error) []byte {
	var payloadBuildingErr error
	var payloadResponse []byte
	if actionErr == nil {
		payloadResponse, payloadBuildingErr = json.Marshal(ActionMessageResponse{
			Status: "success",
		})
	} else {
		payloadResponse, payloadBuildingErr = json.Marshal(ActionMessageResponse{
			Status:  "error",
			Details: actionErr.Error(),
		})
	}
	if payloadBuildingErr != nil {
		panic(payloadBuildingErr)
	}
	return payloadResponse
}

func SendOutputMessage(client *clients.Client, _type string, _payloadResponse []byte) {
	client.Send(OutputMessage{Type: _type, Payload: _payloadResponse})
}

func SendActionMessageResponse(client *clients.Client, _type string, actionErr error) {
	SendOutputMessage(client, _type, BuildActionMessageResponse(actionErr))
}
