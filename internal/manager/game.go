package manager

import (
	"OnlineGame/internal/database"
	"OnlineGame/internal/game"
	"OnlineGame/internal/game/world"
	gamepb "OnlineGame/pkg/pb/go/game"
	"errors"
	"fmt"
)

func (m *Manager) CreateGame(clientID uint, gameName string, worldString string) (*game.Game, error) {
	_currentMatchID, _inGame := m.IsClientInMatch(clientID)
	if _inGame {
		return nil, fmt.Errorf("user currently in match with id=%d", _currentMatchID)
	}

	repo := database.NewMatchRepository(database.GetDB())
	match, err := repo.Create(gameName, worldString, clientID)
	if err != nil {
		return nil, errors.New("failed to create match")
	}

	_game, err := game.NewGame(
		match,
		func() { // On game startup
			fmt.Println("Good luck from Manager!")
			m.broadcastMessage(match.ID, &gamepb.OutputMessage{
				Payload: &gamepb.OutputMessage_Notification{
					Notification: &gamepb.Notification{
						Type:    "match_started",
						Message: "Match started",
						MatchId: uint64(match.ID),
					},
				},
			})
		},
		func() { // On game shut down
			fmt.Println("Bye from Manager!")
			m.broadcastMessage(match.ID, &gamepb.OutputMessage{
				Payload: &gamepb.OutputMessage_Notification{
					Notification: &gamepb.Notification{
						Type:     "match_ended",
						Message:  "Match ended",
						MatchId:  uint64(match.ID),
						PlayerId: uint64(match.WinnerPlayerID),
					},
				},
			})
			database.GetDB().Save(match)
			m.deleteGame(match.ID)
			m.deleteLobby(match.ID)
		},
	)

	if err != nil {
		return nil, err
	}
	m.createEmptyLobby(match.ID)
	m.addIdleGame(_game)
	return _game, nil
}

func (m *Manager) StartGame(clientID uint, matchID uint) error {
	_game, _foundGame := m.idleGames[matchID]
	if !_foundGame {
		return errors.New("game not in idle status or game does not exist")
	}
	if clientID != _game.Match.OwnerPlayerID {
		return errors.New("you are not the owner of the game")
	}

	minPlayers, _, err := world.ParseWorldPropsFromString(_game.Match.WorldString)
	if err != nil {
		return fmt.Errorf("could not parse world props: %w", err)
	}

	_lobby := m.lobby[matchID]
	if len(_lobby) < int(minPlayers) {
		return fmt.Errorf("provided %d players. for start needed minimum %d", len(_lobby), minPlayers)
	}

	err = _game.Start(_lobby)
	if err != nil {
		return fmt.Errorf("error starting game: %w", err)
	}

	m.deleteIdleGame(matchID)
	m.addActiveGame(_game)
	return nil
}

func (m *Manager) deleteGame(matchID uint) *game.Game {
	var _game *game.Game
	_game = m.deleteActiveGame(matchID)
	if _game == nil {
		_game = m.deleteIdleGame(matchID)
	}
	return _game
}
