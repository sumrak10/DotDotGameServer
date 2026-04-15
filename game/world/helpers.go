package world

import "OnlineGame/game/world/nodes"

func (w *World) getNodeNeighbors(nodeID nodes.NodeID) []*nodes.Node {
	neighborIDs := w.nodeNeighbors[nodeID]
	neighbors := make([]*nodes.Node, 0, len(neighborIDs))

	for _, id := range neighborIDs {
		if n, ok := w.Nodes[id]; ok {
			neighbors = append(neighbors, n)
		}
	}

	return neighbors
}

func (w *World) getNodeEdgeByN1N2(n1ID, n2ID nodes.NodeID) *nodes.NodeEdge {
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

func (w *World) GetDelta() float64 {
	return w.delta
}
