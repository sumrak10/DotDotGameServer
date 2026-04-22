package game

import (
	"OnlineGame/internal/clients"
	"OnlineGame/internal/config"
	"OnlineGame/internal/database"
	"OnlineGame/internal/game/world"
	gamepb "OnlineGame/pkg/pb/go/game"
	"context"
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
	inputQueue chan *inputQueueMessage
	running    atomic.Bool
	onStartUp  func()
	onShutDown func()
	//Players
	players map[uint]*Player
	// World
	ExampleWorld *world.World
	world        *world.World
}

func NewGame(match *database.Match, onStartUp, onShutDown func()) (*Game, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Game{
		Match: match,

		ctx:        ctx,
		cancel:     cancel,
		inputQueue: make(chan *inputQueueMessage, 1024),

		onStartUp:  onStartUp,
		onShutDown: onShutDown,

		players: make(map[uint]*Player),
	}, nil
}

func (g *Game) Start(clients []*clients.Client) error {
	fmt.Println("Game starting...")
	if g.running.Load() {
		return errors.New("game is already running")
	}

	_world, err := world.NewWorldFromString(g.Match.WorldString)
	if err != nil {
		return fmt.Errorf("error while creating world: %w", err)
	}
	g.world = _world

	playersIdAndStartPositionsMap := make(map[uint]uint)
	for i, client := range clients {
		g.players[client.User.ID] = NewPlayer(client)
		playersIdAndStartPositionsMap[uint(i)] = client.User.ID
	}
	g.world.Init(playersIdAndStartPositionsMap)

	if !g.world.IsInitialized() {
		panic("game is not initialized")
	}

	g.running.Store(true)
	fmt.Println("World initialised")
	go g.run()
	return nil
}

func (g *Game) Stop() {
	g.cancel()
}

func (g *Game) IsRunning() bool {
	return g.running.Load()
}

func (g *Game) run() {
	ticker := time.NewTicker(time.Second / time.Duration(int64(config.Game().TPS)))
	lastTick := time.Now()

	// On Start up
	if g.onStartUp != nil {
		g.onStartUp()
	}
	fmt.Println("Game started")

	for {
		select {
		case <-g.ctx.Done():
			// On shutdown
			fmt.Println("Game stopping...")

			if g.onShutDown != nil {
				g.onShutDown()
			}

			ticker.Stop()
			g.running.Store(false)

			fmt.Println("Game stopped")
			return
		case now := <-ticker.C:
			// Game cycle
			elapsed := now.Sub(lastTick)
			lastTick = now
			actualTPS := 1.0 / elapsed.Seconds()

			g.processInputs()
			g.tick(actualTPS)
		}
	}
}

func (g *Game) tick(currentTPS float64) {
	playersActivesCounter := g.world.Tick()
	deadPlayers := 0
	var winnerID uint
	for playerID, player := range g.players {
		actives, found := playersActivesCounter[playerID]
		if !found || actives == 0 {
			player.IsAlive = false
			deadPlayers++
		} else {
			winnerID = playerID
		}
	}
	if deadPlayers >= len(g.players)-1 {
		g.Match.WinnerPlayerID = winnerID
		g.Stop()
		return
	}

	tickMessage := &gamepb.OutputMessage{
		Payload: &gamepb.OutputMessage_Tick{
			Tick: &gamepb.Tick{
				MatchId:    uint32(g.Match.ID),
				CurrentTps: float32(currentTPS),
				World:      g.world.ToProto(),
			},
		},
	}
	g.broadcastMessage(tickMessage)
}
