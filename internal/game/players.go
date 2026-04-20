package game

import (
	"OnlineGame/internal/clients"
)

type Player struct {
	Client  *clients.Client
	IsAlive bool
}

func NewPlayer(client *clients.Client) *Player {
	return &Player{
		Client:  client,
		IsAlive: true,
	}
}

func (g *Game) ChangePlayerClient(newClient *clients.Client) {
	player, found := g.players[newClient.User.ID]
	if !found {
		panic("Player not found")
	}
	player.Client = newClient
}
