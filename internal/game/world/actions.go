package world

import (
	"OnlineGame/internal/game/world/nodes"
	nodespb "OnlineGame/pkg/pb/go/game/world/nodes"
	"errors"
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
	nodeEdge := w.getNodeEdgeByN1N2(headingFromNode.ID, headingToNode.ID)
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

func (w *World) UpdateNodeType(playerID uint, nodeID nodes.NodeID, NewType nodespb.NodeType) error {
	node := w.Nodes[nodeID]
	if node.OwnerID != playerID {
		return errors.New("player can't change type not own node")
	}
	if node.Value < nodes.NodeTypePropsMap[NewType].TransformCost {
		return errors.New("new type is not enough to transform cost")
	}
	node.Type = NewType
	node.Value -= nodes.NodeTypePropsMap[NewType].TransformCost
	if node.Shield > nodes.NodeTypePropsMap[NewType].MaxShield {
		node.Shield = nodes.NodeTypePropsMap[NewType].MaxShield
	}
	node.ResetTicks()
	return nil
}

func (w *World) SetAlwaysSendArmy(playerID uint, fromNodeID nodes.NodeID, toNodeID nodes.NodeID, mode bool) error {
	node := w.Nodes[fromNodeID]
	if node.OwnerID != playerID {
		return errors.New("player can't set AlwaysSendArmy not own node")
	}
	if fromNodeID == toNodeID {
		return errors.New("can't send army to same node from which it was originally sent ")
	}
	node.IsAlwaysSendArmy = mode
	node.AlwaysSendArmyToNodeID = toNodeID
	return nil
}
