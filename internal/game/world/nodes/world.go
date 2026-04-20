package nodes

type WorldInterface interface {
	GetNodeByID(id NodeID) *Node
	GetNodeEdgeByID(id NodeEdgeID) *NodeEdge
	SendArmy(playerID uint, headingFromID NodeID, headingToID NodeID, value uint) error
}
