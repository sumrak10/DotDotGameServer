package world

import (
	"OnlineGame/game/world/nodes"
	"errors"
)

func (w *World) SendArmy(playerID uint, headingFrom nodes.NodeID, headingTo nodes.NodeID, value uint) error {
	headingFromNode := w.Nodes[headingFrom]
	if headingFromNode.OwnerID != playerID {
		return errors.New("player can't send army from not own node")
	}
	if headingFromNode.Value < value {
		return errors.New("army value is not enough to send army")
	}
	headingToNode := w.Nodes[headingTo]
	nodeEdge := w.getNodeEdgeByN1N2(headingFrom, headingTo)
	if nodeEdge == nil {
		return errors.New("between this nodes no edge")
	}
	headingFromNode.Value -= value

	_army := &nodes.Army{
		ID:          w.armyIDGenerator(),
		Pos:         0,
		NodeEdge:    nodeEdge,
		HeadingFrom: headingFromNode,
		HeadingTo:   headingToNode,
		Value:       value,
	}

	nodeEdge.Armies = append(nodeEdge.Armies, _army)
	return nil
}

func (w *World) UpdateNodeType(playerID uint, nodeID nodes.NodeID, NewType nodes.NodeType) error {
	node := w.Nodes[nodeID]
	if node.OwnerID != playerID {
		return errors.New("player can't change type not own node")
	}
	if node.Value < NewType.TransformCost() {
		return errors.New("new type is not enough to transform cost")
	}
	node.Type = NewType
	node.Value -= NewType.TransformCost()
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
