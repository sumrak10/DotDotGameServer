package nodes

import (
	"math"
)

type NodeID uint64

const maxNodeValue = uint(math.MaxUint32 / 2)

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
	sendArmyTick    uint
}

func (n *Node) ResetTicks() {
	n.produceTick = 0
	n.shieldRegenTick = 0
}

func (n *Node) Tick(world WorldInterface) {
	// Production
	if n.OwnerID != 0 && n.Value < maxNodeValue {
		if n.produceTick >= n.Type.ProduceSpeed() {
			n.produceTick = 0
			n.Value++
		}
		n.produceTick++
	}

	// Shields
	if n.Type.MaxShield() > 0 && n.Shield < n.Type.MaxShield() {
		if n.shieldRegenTick >= n.Type.ShieldRegenSpeed() {
			n.shieldRegenTick = 0
			n.Shield++
		}
		n.shieldRegenTick++
	}

	if n.IsAlwaysSendArmy && n.Value >= 1 {
		if n.sendArmyTick >= n.Type.ProduceSpeed() {
			err := world.SendArmy(
				n.OwnerID,
				n.ID,
				n.AlwaysSendArmyToNodeID,
				n.Value,
			)
			if err != nil {
				panic(err)
			}
			n.sendArmyTick = 0
		}
		n.sendArmyTick++
	}
}
