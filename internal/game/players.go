package game

import (
	"OnlineGame/internal/clients"
)

type Player struct {
	GID     uint
	Client  *clients.Client
	IsAlive bool
}

func NewPlayer(GID uint, client *clients.Client) *Player {
	return &Player{
		GID:     GID,
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
