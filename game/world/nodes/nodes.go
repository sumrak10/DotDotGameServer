package nodes

type NodeID uint64

type Node struct {
	ID NodeID `json:"id"`

	// Position
	PosX int `json:"pos_x"`
	PosY int `json:"pos_y"`

	// Props
	OwnerID uint     `json:"owner_id"`
	Type    NodeType `json:"type"`
	Shield  uint     `json:"shield"`
	Value   uint     `json:"value"`

	// Behavior
	IsAlwaysSendArmy       bool   `json:"is_always_send_army"`
	AlwaysSendArmyToNodeID NodeID `json:"always_send_army_to_node_id"`

	// Tick counters
	produceTick     uint
	shieldRegenTick uint
}

func (n *Node) UpdateType(newType NodeType) {
	n.Type = newType
	if newType == DefaultNodeType {
		return
	}
	n.produceTick = 0
	n.shieldRegenTick = 0
}

func (n *Node) Tick() {
	// Production
	if n.OwnerID != 0 {
		if n.produceTick >= n.Type.ProduceSpeed() {
			n.produceTick = 0
			n.Value++
		}
		n.produceTick++
	}

	// Shields
	if n.Type.MaxShield() > 0 {
		if n.Shield >= n.Type.MaxShield() {
			n.shieldRegenTick = 0
		} else if n.shieldRegenTick >= n.Type.ShieldRegenSpeed() {
			n.shieldRegenTick = 0
			n.Value++
		}
		n.shieldRegenTick++
	}
}
