package nodes

type NodeEdgeKey struct {
	N1 NodeID
	N2 NodeID
}

func NewNodeEdgeKey(id1, id2 NodeID) NodeEdgeKey {
	if id1 < id2 {
		return NodeEdgeKey{N1: id1, N2: id2}
	}
	return NodeEdgeKey{N1: id2, N2: id1}
}

type NodeEdgeID uint64

type NodeEdge struct {
	ID     NodeEdgeID `json:"id"`
	N1ID   NodeID     `json:"n1_id"`
	N2ID   NodeID     `json:"n2_id"`
	Length uint       `json:"length"`

	Armies []*Army `json:"armies"`
}

func (n *NodeEdge) Tick(world WorldInterface) {
	writeIdx := 0

	for i := 0; i < len(n.Armies); i++ {
		army := n.Armies[i]

		army.Tick(world)

		if army.Value > 0 {
			n.Armies[writeIdx] = army
			writeIdx++
		}
	}

	for j := writeIdx; j < len(n.Armies); j++ {
		n.Armies[j] = nil
	}
	n.Armies = n.Armies[:writeIdx]
}
