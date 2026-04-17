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
	inputQueue chan PlayerInput
	running    atomic.Bool
	OnShutdown func()
	//Players
	players       map[uint]*Player
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

		ctx:        ctx,
		cancel:     cancel,
		inputQueue: make(chan PlayerInput, 1024),

		players:       make(map[uint]*Player),
		ownerPlayerID: ownerPlayerID,

		WorldBuilderName: worldBuilderName,
		ExampleWorld:     world.GetPresetVault().BuildWorld(worldBuilderName),
	}, nil
}

func (g *Game) Start(ownerID uint, clients []*clients.Client) error {
	if g.running.Load() {
		return errors.New("game is already running")
	}
	if ownerID != g.ownerPlayerID {
		return errors.New("you are not the owner of the game")
	}
	if len(clients) < int(g.ExampleWorld.MinPlayers) {
		return errors.New(fmt.Sprintf("provided %d players. for start needed minimum %d", len(g.players), g.ExampleWorld.MinPlayers))
	}
	g.running.Store(true)
	g.world = world.GetPresetVault().BuildWorld(g.WorldBuilderName)

	_indexedPlayers := make(map[uint]uint)
	for i, client := range clients {
		g.players[client.User.ID] = NewPlayer(client)
		_indexedPlayers[uint(i)] = client.User.ID
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

	tps := int64(config.Game().TPS)
	frameDuration := time.Second / time.Duration(tps)
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()
	fixedDelta := 1.0 / float64(tps)
	lastTick := time.Now()
	for {
		select {
		case <-g.ctx.Done():
			return
		case now := <-ticker.C:
			elapsed := now.Sub(lastTick)
			lastTick = now
			actualTPS := 1.0 / elapsed.Seconds()
			g.tick(fixedDelta, actualTPS)
		}
	}
}

func (g *Game) tick(delta float64, currentTPS float64) {
	g.processInputs()
	_playersActivesCounter := g.world.Tick(delta)
	_deadPlayers := 0
	var _winnerID uint
	for playerID, count := range _playersActivesCounter {
		if count == 0 {
			g.players[playerID].IsAlive = false
			_deadPlayers++
		} else {
			_winnerID = playerID
		}
	}
	if _deadPlayers >= len(g.players)-1 {
		g.NotifyGameWinner(_winnerID)
		g.Stop()
		return
	}

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
	for _, player := range g.players {
		if player.Client.Conn != nil {
			player.Client.Send(msg)
		}
	}
}
