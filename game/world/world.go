package world

import (
	"OnlineGame/game/world/nodes"
	"math/rand/v2"
)

type World struct {
	// Props
	MinPlayers        uint8         `json:"min_players"`
	MaxPlayers        uint8         `json:"max_players"`
	PlayersStartNodes []*nodes.Node `json:"players_start_nodes"`

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

	// Physics
	delta float64
}

func NewWorld(minPlayers, maxPlayers uint8) *World {
	return &World{
		MinPlayers:        minPlayers,
		MaxPlayers:        maxPlayers,
		PlayersStartNodes: make([]*nodes.Node, 0, maxPlayers),

		Nodes: make(map[nodes.NodeID]*nodes.Node),

		NodeEdges:       make([]*nodes.NodeEdge, 0),
		nodeEdgeIDMap:   make(map[nodes.NodeEdgeID]*nodes.NodeEdge),
		nodeEdgesKeyMap: make(map[nodes.NodeEdgeKey]*nodes.NodeEdge),
		nodeNeighbors:   make(map[nodes.NodeID][]nodes.NodeID),

		nodeIDGenerator:     nodes.CreateNodeIDGenerator(),
		nodeEdgeIDGenerator: nodes.CreateNodeEdgeIDGenerator(),
		armyIDGenerator:     nodes.CreateArmyIDGenerator(),
	}
}

func (w *World) Init(players map[uint]uint) {
	minValue := uint64(1)
	maxValue := uint64(20)
	// Init other nodes
	for _, node := range w.Nodes {
		node.Value = uint(minValue + rand.Uint64N(maxValue-minValue+1))
	}
	// Init start nodes for players
	for i, playerID := range players {
		playerStartNode := w.PlayersStartNodes[i]
		playerStartNode.OwnerID = playerID
		playerStartNode.Value = uint(10)
		for _, neighborNode := range w.getNodeNeighbors(playerStartNode.ID) {
			neighborNode.Value = uint(10)
		}
	}
}

func (w *World) Tick(delta float64) map[uint]uint {
	w.delta = delta
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
