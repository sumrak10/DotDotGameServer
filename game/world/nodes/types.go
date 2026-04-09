package nodes

type NodeType uint8

const (
	DefaultNodeType NodeType = iota
	ProductionNodeType
	FortressNodeType
)

type NodeTypeProps struct {
	TransformCost    uint
	ProduceSpeed     uint
	MaxShield        uint
	ShieldRegenSpeed uint
}

var nodeTypePropsMap = map[NodeType]NodeTypeProps{
	DefaultNodeType: {
		TransformCost:    0,
		ProduceSpeed:     60,
		MaxShield:        5,
		ShieldRegenSpeed: 30,
	},
	ProductionNodeType: {
		TransformCost:    20,
		ProduceSpeed:     120,
		MaxShield:        0,
		ShieldRegenSpeed: 0,
	},
	FortressNodeType: {
		TransformCost:    20,
		ProduceSpeed:     30,
		MaxShield:        20,
		ShieldRegenSpeed: 60,
	},
}

func (n NodeType) Props() NodeTypeProps {
	return nodeTypePropsMap[n]
}

func (n NodeType) TransformCost() uint    { return n.Props().TransformCost }
func (n NodeType) ProduceSpeed() uint     { return n.Props().ProduceSpeed }
func (n NodeType) MaxShield() uint        { return n.Props().MaxShield }
func (n NodeType) ShieldRegenSpeed() uint { return n.Props().ShieldRegenSpeed }
