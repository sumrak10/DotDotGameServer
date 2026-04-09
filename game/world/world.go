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

	// NodeEdge
	NodeEdges     []*nodes.NodeEdge `json:"node_edges"`
	nodeEdgesMap  map[nodes.NodeEdgeKey]*nodes.NodeEdge
	nodeNeighbors map[nodes.NodeID][]nodes.NodeID

	// UID generators
	nodeIDGenerator     func() nodes.NodeID
	nodeEdgeIDGenerator func() nodes.NodeEdgeID
	armyIDGenerator     func() nodes.ArmyID
}

func NewWorld(minPlayers, maxPlayers uint8) *World {
	return &World{
		MinPlayers:        minPlayers,
		MaxPlayers:        maxPlayers,
		PlayersStartNodes: make([]*nodes.Node, 0, maxPlayers),

		Nodes: make(map[nodes.NodeID]*nodes.Node),

		NodeEdges:     make([]*nodes.NodeEdge, 0),
		nodeEdgesMap:  make(map[nodes.NodeEdgeKey]*nodes.NodeEdge),
		nodeNeighbors: make(map[nodes.NodeID][]nodes.NodeID),

		nodeIDGenerator:     nodes.CreateNodeIDGenerator(),
		nodeEdgeIDGenerator: nodes.CreateNodeEdgeIDGenerator(),
		armyIDGenerator:     nodes.CreateArmyIDGenerator(),
	}
}

func (w *World) Init(players map[uint]uint) {
	minValue := uint64(10)
	maxValue := uint64(20)
	// Init other nodes
	for _, node := range w.Nodes {
		node.Value = uint(minValue + rand.Uint64N(maxValue-minValue+1))
	}
	// Init start nodes for players
	for i, playerID := range players {
		playerStartNode := w.PlayersStartNodes[i]
		playerStartNode.OwnerID = playerID
		playerStartNode.Value = uint(minValue)
		for _, neighborNode := range w.getNodeNeighbors(playerStartNode.ID) {
			neighborNode.Value = uint(minValue)
		}
	}
}

func (w *World) Tick() {
	for _, node := range w.Nodes {
		// Node production, Node shields
		node.Tick()
		// Node AlwaysSendArmyToNodeID logic
		if node.IsAlwaysSendArmy && node.Value >= 1 {
			err := w.SendArmy(
				node.OwnerID,
				node.ID,
				node.AlwaysSendArmyToNodeID,
				node.Value,
			)
			if err != nil {
				panic("its never happened")
			}
		}
	}
	// Delete zero armies from NodeEdge.Armies
	for _, nodeEdge := range w.NodeEdges {
		nodeEdge.Tick()
	}
}
