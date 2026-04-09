package game

import (
	"OnlineGame/clients"
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

	players       map[uint]*clients.Client
	ownerPlayerID uint
	world         *world.World
}

func NewGame(match *database.Match, worldBuilderName string, ownerPlayerID uint) (*Game, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Game{
		Match: match,

		ctx:    ctx,
		cancel: cancel,

		players:       make(map[uint]*clients.Client),
		ownerPlayerID: ownerPlayerID,
		world:         world.GetPresetVault().BuildWorld(worldBuilderName),
	}, nil
}

func (g *Game) Start(ownerID uint, players []*clients.Client) error {
	if g.running.Load() {
		return errors.New("game is already running")
	}
	if ownerID != g.ownerPlayerID {
		return errors.New("you are not the owner of the game")
	}
	if len(players) < int(g.world.MinPlayers) {
		return errors.New(fmt.Sprintf("provided %d players. for start needed minimum %d", len(g.players), g.world.MinPlayers))
	}
	g.running.Store(true)

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

	for {
		select {
		case <-g.ctx.Done():
			return
		default:
			g.tick()
			time.Sleep(time.Millisecond * 16)
		}
	}
}

func (g *Game) tick() {
	g.world.Tick()

	_payload, err := json.Marshal(TickMessagePayload{
		MatchID: g.Match.ID,
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
