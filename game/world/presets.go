package world

import (
	"errors"
	"fmt"
	"sync"
)

func defaultWorldBuilder() *World {
	w := NewWorld(1, 2)
	// Nodes
	n1 := w.addNode(0, 0, 0)
	n2 := w.addNode(0, 1, 0)
	n3 := w.addNode(1, 0, 0)
	n4 := w.addNode(1, 1, 0)
	// Edges
	w.addEdgeNode(n1, n2)
	w.addEdgeNode(n2, n3)
	w.addEdgeNode(n3, n4)
	w.addEdgeNode(n4, n1)
	// Start positions
	w.addPlayerStartNode(n1)
	w.addPlayerStartNode(n3)

	return w
}

var (
	instance *PresetVault
	once     sync.Once
)

type PresetVault struct {
	builders map[string]func() *World
}

func NewPresetVault() *PresetVault {
	return &PresetVault{
		builders: make(map[string]func() *World),
	}
}

func GetPresetVault() *PresetVault {
	once.Do(func() {
		instance = NewPresetVault()
		instance.addBuilder("default", defaultWorldBuilder)
	})
	return instance
}

func (p *PresetVault) addBuilder(name string, builder func() *World) {
	err := p.testBuilder(name, builder)
	if err != nil {
		panic(err)
	}
	p.builders[name] = builder
}

func (p *PresetVault) BuildWorld(name string) *World {
	return p.builders[name]()
}

func (p *PresetVault) testBuilder(name string, builder func() *World) error {
	testWorld := builder()
	if uint8(len(testWorld.PlayersStartNodes)) != testWorld.MaxPlayers {
		return errors.New(fmt.Sprintf("'%s' builder: not all players start nodes specified", name))
	}
	return nil
}
