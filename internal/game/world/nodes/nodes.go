package nodes

import (
	nodespb "OnlineGame/pkg/pb/go/game/world/nodes"
	"math"
)

type NodeID uint64

const maxNodeValue = uint(math.MaxUint32 / 2)

type Node struct {
	ID NodeID

	// Position
	PosX int
	PosY int

	// Props
	OwnerID uint
	Type    nodespb.NodeType
	Shield  uint
	Value   uint

	// Behavior
	IsAlwaysSendArmy       bool
	AlwaysSendArmyToNodeID NodeID

	// Tick counters
	produceTick     uint
	shieldRegenTick uint
	sendArmyTick    uint
}

func (n *Node) ToProto() *nodespb.Node {
	return &nodespb.Node{
		Id:                     uint64(n.ID),
		PosX:                   int32(n.PosX),
		PosY:                   int32(n.PosY),
		OwnerId:                uint64(n.OwnerID),
		Type:                   n.Type,
		Shield:                 uint32(n.Shield),
		Value:                  uint32(n.Value),
		IsAlwaysSendArmy:       n.IsAlwaysSendArmy,
		AlwaysSendArmyToNodeId: uint64(n.AlwaysSendArmyToNodeID),
	}
}

func (n *Node) ResetTicks() {
	n.produceTick = 0
	n.shieldRegenTick = 0
}

func (n *Node) Tick(world WorldInterface) {
	// Production
	if n.OwnerID != 0 && n.Value < maxNodeValue {
		if n.produceTick >= NodeTypePropsMap[n.Type].ProduceSpeed {
			n.produceTick = 0
			n.Value++
		}
		n.produceTick++
	}

	// Shields
	if NodeTypePropsMap[n.Type].MaxShield > 0 && n.Shield < NodeTypePropsMap[n.Type].MaxShield {
		if n.shieldRegenTick >= NodeTypePropsMap[n.Type].ShieldRegenSpeed {
			n.shieldRegenTick = 0
			n.Shield++
		}
		n.shieldRegenTick++
	}

	if n.IsAlwaysSendArmy && n.Value >= 1 {
		if n.sendArmyTick >= NodeTypePropsMap[n.Type].ProduceSpeed {
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
