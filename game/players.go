package game

import (
	"OnlineGame/clients"
	"OnlineGame/utils"
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

func (g *Game) OwnerPlayerID() uint {
	return g.ownerPlayerID
}

func (g *Game) IsPlayerAdmin(playerID uint) bool {
	return playerID == g.ownerPlayerID
}

func (g *Game) MakeAdminRandomPlayer() {
	playersIDs := utils.MapKeys(g.players)
	g.ownerPlayerID = utils.ListRandomElement(playersIDs)
}

func (g *Game) ChangePlayerClient(newClient *clients.Client) {
	player, found := g.players[newClient.User.ID]
	if !found {
		panic("Player not found")
	}
	player.Client = newClient
}
