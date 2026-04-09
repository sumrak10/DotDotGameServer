package nodes

import "sync/atomic"

func CreateNodeIDGenerator() func() NodeID {
	var counter uint64

	return func() NodeID { // 2. Анонимная функция не может иметь имени GenerateNodeID
		return NodeID(atomic.AddUint64(&counter, 1))
	}
}

func CreateNodeEdgeIDGenerator() func() NodeEdgeID {
	var counter uint64

	return func() NodeEdgeID {
		return NodeEdgeID(atomic.AddUint64(&counter, 1))
	}
}

func CreateArmyIDGenerator() func() ArmyID {
	var counter uint64

	return func() ArmyID {
		return ArmyID(atomic.AddUint64(&counter, 1))
	}
}
