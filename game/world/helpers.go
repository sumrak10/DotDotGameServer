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
	return w.nodeEdgesMap[nodes.NewNodeEdgeKey(n1ID, n2ID)]
}
