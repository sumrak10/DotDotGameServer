package world

import (
	"OnlineGame/internal/config"
	"OnlineGame/internal/game/world/nodes"
	nodespb "OnlineGame/pkg/pb/go/game/world/nodes"
	"errors"
	"fmt"
	"math"
)

// fetching components utils

func (w *World) GetNodeNeighbors(nodeID nodes.NodeID) []*nodes.Node {
	neighborIDs := w.nodeNeighbors[nodeID]
	neighbors := make([]*nodes.Node, 0, len(neighborIDs))

	for _, id := range neighborIDs {
		if n, ok := w.Nodes[id]; ok {
			neighbors = append(neighbors, n)
		}
	}

	return neighbors
}

func (w *World) GetNodeEdgeByN1N2(n1ID, n2ID nodes.NodeID) *nodes.NodeEdge {
	nodeEdge, found := w.nodeEdgesKeyMap[nodes.NewNodeEdgeKey(n1ID, n2ID)]
	if !found {
		return nil
	}
	return nodeEdge
}

func (w *World) GetNodeByID(id nodes.NodeID) *nodes.Node {
	node, found := w.Nodes[id]
	if !found {
		return nil
	}
	return node
}

func (w *World) GetNodeEdgeByID(id nodes.NodeEdgeID) *nodes.NodeEdge {
	nodeEdge, found := w.nodeEdgeIDMap[id]
	if !found {
		return nil
	}
	return nodeEdge
}

// adding components utils

func (w *World) addNode(x, y int, ownerID uint) (*nodes.Node, error) {
	if w.isInitialized {
		return nil, errors.New("can't add node to initialized world")
	}
	id := w.nodeIDGenerator()
	node := &nodes.Node{
		ID: id,

		// Position
		PosX: x,
		PosY: y,

		// Props
		OwnerID: ownerID,
		Type:    nodespb.NodeType_DefaultNodeType,
		Shield:  nodes.NodeTypePropsMap[nodespb.NodeType_DefaultNodeType].MaxShield,
		Value:   10,
	}
	w.Nodes[id] = node
	return node, nil
}

func (w *World) addEdgeNode(n1ID nodes.NodeID, n2ID nodes.NodeID) (*nodes.NodeEdge, error) {
	if w.isInitialized {
		return nil, errors.New("can't add node edge to initialized world")
	}
	n1, n1Found := w.Nodes[n1ID]
	if !n1Found {
		return nil, fmt.Errorf("node#%d not found", uint(n1ID))
	}
	n2, n2Found := w.Nodes[n2ID]
	if !n2Found {
		return nil, fmt.Errorf("node#%d not found", uint(n2ID))
	}
	length := math.Abs(
		math.Sqrt(
			math.Pow(float64(n2.PosX-n1.PosX), 2)+
				math.Pow(float64(n2.PosY-n1.PosY), 2),
		),
	) * float64(config.Game().ValuesScaleCoef)
	id := w.nodeEdgeIDGenerator()
	nodeEdge := &nodes.NodeEdge{
		ID:     id,
		N1ID:   n1ID,
		N2ID:   n2ID,
		Length: uint(length),
		Armies: make([]*nodes.Army, 0),
	}
	w.NodeEdges = append(w.NodeEdges, nodeEdge)
	w.nodeEdgeIDMap[nodeEdge.ID] = nodeEdge
	w.nodeEdgesKeyMap[nodes.NewNodeEdgeKey(n1.ID, n2.ID)] = nodeEdge
	w.nodeNeighbors[n1.ID] = append(w.nodeNeighbors[n1.ID], n2.ID)
	w.nodeNeighbors[n2.ID] = append(w.nodeNeighbors[n2.ID], n1.ID)
	return nodeEdge, nil
}

func (w *World) addPlayerStartNode(nodeID nodes.NodeID) error {
	if w.isInitialized {
		return errors.New("can't player start position to initialized world")
	}
	node, found := w.Nodes[nodeID]
	if !found {
		return fmt.Errorf("node#%d not found", uint(nodeID))
	}
	w.PlayersStartNodes = append(w.PlayersStartNodes, node)
	return nil
}
