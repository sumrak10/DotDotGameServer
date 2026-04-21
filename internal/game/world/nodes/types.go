package nodes

import nodespb "OnlineGame/pkg/pb/go/game/world/nodes"

type NodeTypeProps struct {
	TransformCost    uint
	ProduceSpeed     uint
	MaxShield        uint
	ShieldRegenSpeed uint
}

var NodeTypePropsMap = map[nodespb.NodeType]NodeTypeProps{
	nodespb.NodeType_DefaultNodeType: {
		TransformCost:    0, // can't transform to this type
		ProduceSpeed:     30,
		MaxShield:        5,
		ShieldRegenSpeed: 30,
	},
	nodespb.NodeType_ProductionNodeType: {
		TransformCost:    10,
		ProduceSpeed:     15,
		MaxShield:        0,
		ShieldRegenSpeed: 0,
	},
	nodespb.NodeType_FortressNodeType: {
		TransformCost:    10,
		ProduceSpeed:     120,
		MaxShield:        20,
		ShieldRegenSpeed: 15,
	},
}
