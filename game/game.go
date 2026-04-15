package game

import (
	"OnlineGame/clients"
	"OnlineGame/config"
	"OnlineGame/database"
	"OnlineGame/game/world"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type Game struct {
	Match *database.Match
	// Context
	ctx        context.Context
	cancel     context.CancelFunc
	running    atomic.Bool
	OnShutdown func()
	//Players
	players       map[uint]*clients.Client
	ownerPlayerID uint
	// World
	WorldBuilderName string
	ExampleWorld     *world.World
	world            *world.World
}

func NewGame(match *database.Match, worldBuilderName string, ownerPlayerID uint) (*Game, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Game{
		Match: match,

		ctx:    ctx,
		cancel: cancel,

		players:       make(map[uint]*clients.Client),
		ownerPlayerID: ownerPlayerID,

		WorldBuilderName: worldBuilderName,
		ExampleWorld:     world.GetPresetVault().BuildWorld(worldBuilderName),
	}, nil
}

func (g *Game) Start(ownerID uint, players []*clients.Client) error {
	if g.running.Load() {
		return errors.New("game is already running")
	}
	if ownerID != g.ownerPlayerID {
		return errors.New("you are not the owner of the game")
	}
	if len(players) < int(g.ExampleWorld.MinPlayers) {
		return errors.New(fmt.Sprintf("provided %d players. for start needed minimum %d", len(g.players), g.ExampleWorld.MinPlayers))
	}
	g.running.Store(true)
	g.world = world.GetPresetVault().BuildWorld(g.WorldBuilderName)

	_indexedPlayers := make(map[uint]uint)
	for i, player := range players {
		g.players[player.User.ID] = player
		_indexedPlayers[uint(i)] = player.User.ID
	}
	g.world.Init(_indexedPlayers)
	g.NotifyGameStarted()
	go g.run()
	return nil
}

func (g *Game) Stop() {
	g.OnShutdown()
	g.NotifyGameEnded()
	g.running.Store(false)
	g.cancel()
}

func (g *Game) IsRunning() bool {
	return g.running.Load()
}

func (g *Game) run() {
	defer g.Stop()

	tickDuration := time.Second / time.Duration(config.Game().TPS)
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	lastTick := time.Now()
	for range ticker.C {
		select {
		case <-g.ctx.Done():
			return
		default:
			now := time.Now()
			elapsed := now.Sub(lastTick)
			lastTick = now
			actualTPS := 1.0 / elapsed.Seconds()
			g.tick(elapsed.Seconds(), actualTPS)
		}
	}
}

func (g *Game) tick(delta float64, currentTPS float64) {
	g.world.Tick(delta)

	_payload, err := json.Marshal(TickMessagePayload{
		MatchID:    g.Match.ID,
		CurrentTPS: float32(currentTPS),
		World: WorldTickMessagePayload{
			Nodes:     g.world.Nodes,
			NodeEdges: g.world.NodeEdges,
		},
	})
	if err != nil {
		panic(err)
	}
	msg := OutputMessage{
		Type:    "match_tick",
		Payload: _payload,
	}
	g.broadcast(msg)
}

func (g *Game) broadcast(msg any) {
	for _, client := range g.players {
		client.Send(msg)
	}
}
