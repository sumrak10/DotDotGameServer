package world

import (
	"OnlineGame/internal/config"
	"OnlineGame/internal/game/world/nodes"
	nodespb "OnlineGame/pkg/pb/go/game/world/nodes"
	"math"
)

func (w *World) addNode(x, y int, ownerID uint) *nodes.Node {
	_id := w.nodeIDGenerator()
	_node := &nodes.Node{
		ID: _id,

		// Position
		PosX: x,
		PosY: y,

		// Props
		OwnerID: ownerID,
		Type:    nodespb.NodeType_DefaultNodeType,
		Shield:  nodes.NodeTypePropsMap[nodespb.NodeType_DefaultNodeType].MaxShield,
		Value:   10,
	}
	w.Nodes[_id] = _node
	return _node
}

func (w *World) addEdgeNode(n1 *nodes.Node, n2 *nodes.Node) *nodes.NodeEdge {
	length := math.Abs(
		math.Sqrt(
			math.Pow(float64(n2.PosX-n1.PosX), 2) +
				math.Pow(float64(n2.PosY-n1.PosY), 2),
		),
	)
	length *= float64(config.Game().MotionScale)
	_id := w.nodeEdgeIDGenerator()
	_nodeEdge := &nodes.NodeEdge{
		ID:     _id,
		N1ID:   n1.ID,
		N2ID:   n2.ID,
		Length: uint(length),
		Armies: make([]*nodes.Army, 0),
	}
	w.NodeEdges = append(w.NodeEdges, _nodeEdge)
	w.nodeEdgeIDMap[_nodeEdge.ID] = _nodeEdge
	w.nodeEdgesKeyMap[nodes.NewNodeEdgeKey(n1.ID, n2.ID)] = _nodeEdge
	w.nodeNeighbors[n1.ID] = append(w.nodeNeighbors[n1.ID], n2.ID)
	w.nodeNeighbors[n2.ID] = append(w.nodeNeighbors[n2.ID], n1.ID)
	return _nodeEdge
}

func (w *World) addPlayerStartNode(node *nodes.Node) {
	w.PlayersStartNodes = append(w.PlayersStartNodes, node)
}
