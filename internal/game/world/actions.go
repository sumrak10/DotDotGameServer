package world

import (
	"OnlineGame/internal/game/world/nodes"
	nodespb "OnlineGame/pkg/pb/go/game/world/nodes"
	"errors"
	"fmt"
)

func (w *World) SendArmy(playerID uint, headingFromID nodes.NodeID, headingToID nodes.NodeID, value uint) error {
	headingFromNode, _found := w.Nodes[headingFromID]
	if !_found {
		return errors.New("heading from node not found")
	}
	if headingFromNode.OwnerID != playerID {
		return errors.New("player can't send army from not own node")
	}
	if headingFromNode.Value < value {
		return errors.New("army value is not enough to send army")
	}
	headingToNode, _found := w.Nodes[headingToID]
	if !_found {
		return errors.New("heading to node not found")
	}
	nodeEdge := w.GetNodeEdgeByN1N2(headingFromNode.ID, headingToNode.ID)
	if nodeEdge == nil {
		return errors.New("between this nodes no edge")
	}
	headingFromNode.Value -= value

	_army := &nodes.Army{
		ID:            w.armyIDGenerator(),
		Pos:           0,
		NodeEdgeID:    nodeEdge.ID,
		HeadingFromID: headingFromNode.ID,
		HeadingToID:   headingToNode.ID,
		OwnerPlayerID: playerID,
		Value:         value,
	}

	nodeEdge.Armies = append(nodeEdge.Armies, _army)
	return nil
}

func (w *World) UpdateNodeType(playerID uint, nodeID nodes.NodeID, newType nodespb.NodeType) error {
	node := w.Nodes[nodeID]
	if node == nil {
		return fmt.Errorf("node#%d not found", nodeID)
	}
	if node.OwnerID != playerID {
		return errors.New("player can't change type not own node")
	}
	nodeTypeProps, found := nodes.NodeTypePropsMap[newType]
	if !found {
		return fmt.Errorf("for type %d not supported", newType)
	}
	transformCost := nodeTypeProps.TransformCost
	if transformCost == 0 {
		return errors.New("cant transform to this type")
	}
	if node.Value < transformCost {
		return errors.New("new type is not enough to transform cost")
	}

	node.Value -= transformCost
	if node.Shield > nodeTypeProps.MaxShield {
		node.Shield = nodeTypeProps.MaxShield
	}
	node.Type = newType
	node.ResetTicks()

	return nil
}

func (w *World) SetAlwaysSendArmy(playerID uint, fromNodeID nodes.NodeID, toNodeID nodes.NodeID, mode bool) error {
	fromNode := w.Nodes[fromNodeID]
	if fromNode == nil {
		return fmt.Errorf("node#%d not found", fromNodeID)
	}
	toNode := w.Nodes[toNodeID]
	if toNode == nil {
		return fmt.Errorf("node#%d not found", toNodeID)
	}
	if fromNode.OwnerID != playerID {
		return errors.New("player can't set AlwaysSendArmy not own node")
	}
	if fromNodeID == toNodeID {
		return errors.New("can't send army to same node from which it was originally sent ")
	}
	fromNode.IsAlwaysSendArmy = mode
	fromNode.AlwaysSendArmyToNodeID = toNode.ID
	return nil
}
