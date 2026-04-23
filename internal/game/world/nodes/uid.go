package nodes

import "sync/atomic"

func CreateNodeIDGenerator() func() NodeID {
	var counter uint32
	return func() NodeID {
		return NodeID(atomic.AddUint32(&counter, 1))
	}
}

func CreateNodeEdgeIDGenerator() func() NodeEdgeID {
	var counter uint32
	return func() NodeEdgeID {
		return NodeEdgeID(atomic.AddUint32(&counter, 1))
	}
}

func CreateArmyIDGenerator() func() ArmyID {
	var counter uint32
	return func() ArmyID {
		return ArmyID(atomic.AddUint32(&counter, 1))
	}
}
