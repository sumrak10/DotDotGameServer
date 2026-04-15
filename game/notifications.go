package game

type NotificationMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`

	MatchID        uint   `json:"match_id"`
	PlayerID       uint   `json:"player_id"`
	PlayerUserName string `json:"player_username"`
}

func (g *Game) NotifyPlayerDisconnnected(disconnectedPlayerID uint) {
	g.broadcast(NotificationMessage{
		Type:     "player_disconnected",
		Message:  "Player disconnected",
		MatchID:  g.Match.ID,
		PlayerID: disconnectedPlayerID,
	})
}

func (g *Game) NotifyPlayerReconnnected(reConnectedPlayerID uint) {
	g.broadcast(NotificationMessage{
		Type:     "player_reconnected",
		Message:  "Player reconnected",
		MatchID:  g.Match.ID,
		PlayerID: reConnectedPlayerID,
	})
}

func (g *Game) NotifyGameStarted() {
	g.broadcast(NotificationMessage{
		Type:    "match_started",
		Message: "Match started",
		MatchID: g.Match.ID,
	})
}

func (g *Game) NotifyGameEnded() {
	g.broadcast(NotificationMessage{
		Type:    "match_ended",
		Message: "Match ended",
		MatchID: g.Match.ID,
	})
}
