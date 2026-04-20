package nodes

import nodespb "OnlineGame/pkg/pb/go/game/world/nodes"

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

func (n *NodeEdge) ToProto() *nodespb.NodeEdge {
	protoArmies := make([]*nodespb.Army, len(n.Armies))
	for _, a := range n.Armies {
		protoArmies = append(protoArmies, a.ToProto())
	}
	return &nodespb.NodeEdge{
		Id:     uint64(n.ID),
		N1Id:   uint64(n.N1ID),
		N2Id:   uint64(n.N2ID),
		Length: uint32(n.Length),
		Armies: protoArmies,
	}
}

func (n *NodeEdge) Tick(world WorldInterface, playerActiveCounter map[uint]uint) {
	N1 := world.GetNodeByID(n.N1ID)
	N2 := world.GetNodeByID(n.N2ID)

	writeIdx := 0
	for i := 0; i < len(n.Armies); i++ {
		army := n.Armies[i]

		if army.HeadingFromID == n.N1ID {
			army.Tick(world, n, N1, N2)
		} else {
			army.Tick(world, n, N2, N1)
		}

		if army.Value > 0 {
			n.Armies[writeIdx] = army
			writeIdx++
			playerActiveCounter[army.OwnerPlayerID]++
		}
	}

	for j := writeIdx; j < len(n.Armies); j++ {
		n.Armies[j] = nil
	}
	n.Armies = n.Armies[:writeIdx]
}
