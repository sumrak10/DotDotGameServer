package world

import (
	"OnlineGame/internal/game/world/nodes"
	worldpb "OnlineGame/pkg/pb/go/game/world"
	nodespb "OnlineGame/pkg/pb/go/game/world/nodes"
	"math/rand/v2"
)

type World struct {
	// Flags
	isInitialized bool

	// Props
	MinPlayers        uint8
	MaxPlayers        uint8
	PlayerIDnGIDMap   map[uint]uint
	PlayersStartNodes map[uint]nodes.NodeID

	// Node
	Nodes map[nodes.NodeID]*nodes.Node

	// NodeEdgeID
	NodeEdges       []*nodes.NodeEdge
	nodeEdgeIDMap   map[nodes.NodeEdgeID]*nodes.NodeEdge
	nodeEdgesKeyMap map[nodes.NodeEdgeKey]*nodes.NodeEdge
	nodeNeighbors   map[nodes.NodeID][]nodes.NodeID

	// UID generators
	nodeIDGenerator     func() nodes.NodeID
	nodeEdgeIDGenerator func() nodes.NodeEdgeID
	armyIDGenerator     func() nodes.ArmyID
}

func (w *World) ToProto() *worldpb.World {
	protoPlayerIDnGIDMap := make(map[uint32]uint32)
	for playerID, playerGID := range w.PlayerIDnGIDMap {
		protoPlayerIDnGIDMap[uint32(playerID)] = uint32(playerGID)
	}
	protoNodes := make(map[uint32]*nodespb.Node)
	for id, n := range w.Nodes {
		protoNodes[uint32(id)] = n.ToProto()
	}
	protoNodeEdges := make([]*nodespb.NodeEdge, 0, len(w.NodeEdges))
	for _, n := range w.NodeEdges {
		protoNodeEdges = append(protoNodeEdges, n.ToProto())
	}
	return &worldpb.World{
		MinPlayers:      uint32(w.MinPlayers),
		MaxPlayers:      uint32(w.MaxPlayers),
		PlayerIdNGidMap: protoPlayerIDnGIDMap,
		Nodes:           protoNodes,
		NodeEdges:       protoNodeEdges,
	}
}

func (w *World) Init(playerGIDnIDMap map[uint]uint) {
	minValue := uint64(1)
	maxValue := uint64(20)
	// Init other nodes
	for _, node := range w.Nodes {
		node.Value = uint(minValue + rand.Uint64N(maxValue-minValue+1))
	}
	// Init start nodes for players
	for playerGID, playerID := range playerGIDnIDMap {
		w.PlayerIDnGIDMap[playerID] = playerGID

		playerStartNodeID := w.PlayersStartNodes[playerGID]
		playerStartNode := w.Nodes[playerStartNodeID]
		playerStartNode.OwnerID = playerGID
		playerStartNode.Value = uint(10)
		for _, neighborNode := range w.GetNodeNeighbors(playerStartNodeID) {
			neighborNode.Value = uint(10)
		}
	}
	w.isInitialized = true
}

func (w *World) Tick() map[uint]uint {
	playersActivesCount := make(map[uint]uint)

	for _, node := range w.Nodes {
		// Node production, Node shields
		node.Tick(w)
		// Node AlwaysSendArmyToNodeID logic
		playersActivesCount[node.OwnerID]++
	}
	// Node edges and armies ticks
	for _, nodeEdge := range w.NodeEdges {
		nodeEdge.Tick(w, playersActivesCount)
	}

	return playersActivesCount
}

func (w *World) IsInitialized() bool {
	return w.isInitialized
}
