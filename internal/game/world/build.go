package world

import (
	"OnlineGame/internal/game/world/nodes"
	"fmt"
	"strconv"
	"strings"
)

var DefaultPresets = map[string]string{
	"default": "2:2;0,0:0,1:1,0:1,1;1,2:2,3:3,4:4,1;1:3",
}

func NewWorldFromString(worldString string) (*World, error) {
	stringSections := strings.Split(worldString, ";")
	propsSection, nodesSection, nodeEdgesSection, playersStartNodesSection := stringSections[0], stringSections[1], stringSections[2], stringSections[3]

	// propsSection
	minPlayers, maxPlayers, err := parsePropsFromPropsSection(propsSection)
	if err != nil {
		return nil, err
	}

	// world creating
	resultWorld := &World{
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

	// nodesSection
	err = resultWorld.newNodesFromString(nodesSection)
	if err != nil {
		return nil, fmt.Errorf("error creating nodes: %w", err)
	}

	// nodeEdgesSection
	err = resultWorld.newNodeEdgesFromString(nodeEdgesSection)
	if err != nil {
		return nil, fmt.Errorf("error creating node edges: %w", err)
	}

	// playersStartNodesSection
	err = resultWorld.newPlayersStartPositionsFromString(playersStartNodesSection)
	if err != nil {
		return nil, fmt.Errorf("error creating players start pos: %v", err)
	}

	if len(resultWorld.PlayersStartNodes) < int(maxPlayers) {
		return nil, fmt.Errorf("players start positions less than max players: %v", err)
	}

	return resultWorld, nil
}

func ParseWorldPropsFromString(worldString string) (uint8, uint8, error) {
	stringSections := strings.Split(worldString, ";")
	return parsePropsFromPropsSection(stringSections[0])
}

func parsePropsFromPropsSection(propsSection string) (uint8, uint8, error) {
	propsSectionParts := strings.Split(propsSection, ":")
	minPlayers, err := strconv.ParseUint(propsSectionParts[0], 10, 8)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing min players: %w from original string %s", err, propsSectionParts[0])
	}
	maxPlayers, err := strconv.ParseUint(propsSectionParts[1], 10, 8)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing max players: %w from original string %s", err, propsSectionParts[0])
	}
	return uint8(minPlayers), uint8(maxPlayers), err
}

func (w *World) newNodesFromString(nodesSection string) error {
	nodesSectionParts := strings.Split(nodesSection, ":")
	for _, part := range nodesSectionParts {
		nodePositions := strings.Split(part, ",")
		if len(nodePositions) != 2 {
			return fmt.Errorf("node positions must be 2. real string: %s", part)
		}
		posX, err := strconv.ParseUint(nodePositions[0], 10, 32)
		if err != nil {
			return err
		}
		posY, err := strconv.ParseUint(nodePositions[1], 10, 32)
		if err != nil {
			return err
		}
		_, err = w.addNode(int(posX), int(posY), 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) newNodeEdgesFromString(nodeEdgesSection string) error {
	edgesSectionParts := strings.Split(nodeEdgesSection, ":")
	for _, part := range edgesSectionParts {
		nodeIDs := strings.Split(part, ",")
		if len(nodeIDs) != 2 {
			return fmt.Errorf("invalid node id: %s", part)
		}
		n1ID, err := strconv.ParseUint(nodeIDs[0], 10, 32)
		if err != nil {
			return err
		}
		n2ID, err := strconv.ParseUint(nodeIDs[1], 10, 32)
		if err != nil {
			return err
		}
		_, err = w.addEdgeNode(nodes.NodeID(n1ID), nodes.NodeID(n2ID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) newPlayersStartPositionsFromString(playersStartNodesSection string) error {
	playersStartNodesSectionParts := strings.Split(playersStartNodesSection, ":")
	for _, nodeIDString := range playersStartNodesSectionParts {
		nodeID, err := strconv.ParseUint(nodeIDString, 10, 32)
		if err != nil {
			return err
		}
		err = w.addPlayerStartNode(nodes.NodeID(nodeID))
		if err != nil {
			return err
		}
	}
	return nil
}
