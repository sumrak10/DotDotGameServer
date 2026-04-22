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
	PlayersStartNodes []*nodes.Node

	// Node
	Nodes map[nodes.NodeID]*nodes.Node `json:"nodes"`

	// NodeEdgeID
	NodeEdges       []*nodes.NodeEdge `json:"node_edges"`
	nodeEdgeIDMap   map[nodes.NodeEdgeID]*nodes.NodeEdge
	nodeEdgesKeyMap map[nodes.NodeEdgeKey]*nodes.NodeEdge
	nodeNeighbors   map[nodes.NodeID][]nodes.NodeID

	// UID generators
	nodeIDGenerator     func() nodes.NodeID
	nodeEdgeIDGenerator func() nodes.NodeEdgeID
	armyIDGenerator     func() nodes.ArmyID
}

func (w *World) ToProto() *worldpb.World {
	protoPlayersStartNodes := make([]*nodespb.Node, 0, w.MaxPlayers)
	for _, n := range w.PlayersStartNodes {
		protoPlayersStartNodes = append(protoPlayersStartNodes, n.ToProto())
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
		MinPlayers:        uint32(w.MinPlayers),
		MaxPlayers:        uint32(w.MaxPlayers),
		PlayersStartNodes: protoPlayersStartNodes,
		Nodes:             protoNodes,
		NodeEdges:         protoNodeEdges,
	}
}

func (w *World) Init(playersIdAndStartPositionsMap map[uint]uint) {
	minValue := uint64(1)
	maxValue := uint64(20)
	// Init other nodes
	for _, node := range w.Nodes {
		node.Value = uint(minValue + rand.Uint64N(maxValue-minValue+1))
	}
	// Init start nodes for players
	for startPosition, playerID := range playersIdAndStartPositionsMap {
		playerStartNode := w.PlayersStartNodes[startPosition]
		playerStartNode.OwnerID = playerID
		playerStartNode.Value = uint(10)
		for _, neighborNode := range w.GetNodeNeighbors(playerStartNode.ID) {
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
