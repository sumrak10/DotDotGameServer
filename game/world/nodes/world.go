package nodes

type WorldInterface interface {
	GetDelta() float64
	GetNodeByID(id NodeID) *Node
	GetNodeEdgeByID(id NodeEdgeID) *NodeEdge
}
