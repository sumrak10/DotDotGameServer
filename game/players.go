package game

import (
	"OnlineGame/utils"
)

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
