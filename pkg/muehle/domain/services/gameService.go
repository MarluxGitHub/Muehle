package services

import (
	"errors"
	"marluxgithub/muehle/pkg/muehle/domain/entities"

	"github.com/google/uuid"
)

type GameService struct {
	Game         entities.Game
	BoardService *BoardService
}

func (gameService *GameService) CreateGame() {
	gameService.Game = entities.Game{
		ID:                 0,
		Board:              gameService.BoardService.CreateBoard(),
		Players:            []entities.Player{},
		CurrentPlayerIndex: 0,
		State:              entities.GameStateWaitingForPlayers,
	}
}

func (gameService *GameService) CreateFakeGame() {
	gameService.Game = entities.Game{
		ID:                 0,
		Board:              gameService.BoardService.CreateBoard(),
		Players:            []entities.Player{},
		CurrentPlayerIndex: 0,
		State:              entities.GameStateMovingStone,
	}

	gameService.Game.Players = append(gameService.Game.Players, entities.Player{
		ID:           0,
		Name:         "Player 1",
		Secret:       "1",
		Color:        entities.ColorWhite,
		Stones:       9,
		PuttedStones: 0,
		IsJumping:    false,
	})

	gameService.Game.Players = append(gameService.Game.Players, entities.Player{
		ID:           1,
		Name:         "Player 2",
		Secret:       "2",
		Color:        entities.ColorBlack,
		Stones:       3,
		PuttedStones: 0,
		IsJumping:    false,
	})

	gameService.Game.Board.Fields[0].Color = entities.ColorWhite
	gameService.Game.Board.Fields[1].Color = entities.ColorWhite
	gameService.Game.Board.Fields[2].Color = entities.ColorWhite
	gameService.Game.Board.Fields[3].Color = entities.ColorWhite
	gameService.Game.Board.Fields[4].Color = entities.ColorWhite
	gameService.Game.Board.Fields[5].Color = entities.ColorWhite
	gameService.Game.Board.Fields[6].Color = entities.ColorWhite
	gameService.Game.Board.Fields[7].Color = entities.ColorWhite
	gameService.Game.Board.Fields[12].Color = entities.ColorWhite

	gameService.Game.Board.Fields[9].Color = entities.ColorBlack
	gameService.Game.Board.Fields[10].Color = entities.ColorBlack
	gameService.Game.Board.Fields[11].Color = entities.ColorBlack
}

func (gameService *GameService) AddPlayer(player entities.Player) (string, error) {
	if len(gameService.Game.Players) >= 2 {
		return "", errors.New("game is full")
	}

	secretCode := genSecretCode()
	player.Secret = secretCode

	if len(gameService.Game.Players) == 0 {
		player.Color = entities.ColorWhite
	} else {
		player.Color = entities.ColorBlack
	}

	gameService.Game.Players = append(gameService.Game.Players, player)

	if len(gameService.Game.Players) == 2 {
		gameService.Game.State = entities.GameStatePuttingStone
	}

	return secretCode, nil
}

func (gameService *GameService) GetGameState() entities.GameState {
	return gameService.Game.State
}

func (gameService *GameService) GetCurrentPlayerColor() entities.Color {
	return gameService.Game.Players[gameService.Game.CurrentPlayerIndex].Color
}

func (gameService *GameService) MovePutStone(fieldIndex int, secretCode string) error {
	if gameService.Game.State != entities.GameStatePuttingStone {
		return errors.New("game is not in putting stone state")
	}

	playerIndex := gameService.getPlayerIndexBySecretCode(secretCode)

	if playerIndex == -1 {
		return errors.New("invalid secret code")
	}

	if gameService.Game.Players[playerIndex].Color != gameService.GetCurrentPlayerColor() {
		return errors.New("player is not the current player")
	}

	if gameService.Game.Board.Fields[fieldIndex].Color != entities.ColorUnknown {
		return errors.New("field is not empty")
	}

	gameService.Game.Board.Fields[fieldIndex].Color = gameService.Game.Players[playerIndex].Color
	if gameService.BoardService.HasPlayerThreeStones(gameService.Game.Board, fieldIndex, gameService.Game.Players[playerIndex].Color) {
		gameService.Game.StateBeforeRemovingStone = entities.GameStatePuttingStone
		gameService.Game.State = entities.GameStateRemovingStone

		return nil
	}

	gameService.Game.Players[playerIndex].Stones++
	gameService.Game.Players[playerIndex].PuttedStones++

	// both players put 9 stones
	if gameService.Game.Players[0].PuttedStones == 9 && gameService.Game.Players[1].PuttedStones == 9 {
		gameService.Game.State = entities.GameStateMovingStone

		gameService.canCurrentPlayerMove()
	}

	gameService.nextPlayer()

	return nil
}

func (gameService *GameService) MoveStone(fromFieldIndex, toFieldIndex int, secretCode string) error {

	playerIndex := gameService.getPlayerIndexBySecretCode(secretCode)

	if playerIndex == -1 {
		return errors.New("invalid secret code")
	}

	if gameService.Game.Players[playerIndex].Color != gameService.GetCurrentPlayerColor() {
		return errors.New("player is not the current player")
	}

	if gameService.Game.State != entities.GameStateMovingStone {
		return errors.New("game is not in moving stone state")
	}

	if gameService.Game.Board.Fields[fromFieldIndex].Color != gameService.Game.Players[playerIndex].Color {
		return errors.New("field is not owned by player")
	}

	if gameService.Game.Board.Fields[toFieldIndex].Color != entities.ColorUnknown {
		return errors.New("field is not empty")
	}

	isfineMove := gameService.Game.Players[playerIndex].IsJumping

	// enthält das fromField die Nachbars das toField
	for _, neighbor := range gameService.Game.Board.Fields[fromFieldIndex].Neighbors {
		if neighbor.Index == toFieldIndex {
			isfineMove = true
			break
		}
	}

	if isfineMove == false {
		return errors.New("invalid move")
	}

	gameService.Game.Board.Fields[toFieldIndex].Color = gameService.Game.Players[playerIndex].Color
	gameService.Game.Board.Fields[fromFieldIndex].Color = entities.ColorUnknown

	if gameService.BoardService.HasPlayerThreeStones(gameService.Game.Board, toFieldIndex, gameService.Game.Players[playerIndex].Color) {
		gameService.Game.StateBeforeRemovingStone = entities.GameStateMovingStone
		gameService.Game.State = entities.GameStateRemovingStone

		return nil
	}

	gameService.nextPlayer()

	return nil
}

func (gameService *GameService) nextPlayer() {
	gameService.Game.CurrentPlayerIndex = (gameService.Game.CurrentPlayerIndex + 1) % len(gameService.Game.Players)
}

func (gameService *GameService) RemoveStone(fieldIndex int, secretCode string) error {
	if gameService.Game.State != entities.GameStateRemovingStone {
		return errors.New("game is not in removing stone state")
	}

	playerIndex := gameService.getPlayerIndexBySecretCode(secretCode)

	if playerIndex == -1 {
		return errors.New("invalid secret code")
	}

	if gameService.Game.Players[playerIndex].Color != gameService.GetCurrentPlayerColor() {
		return errors.New("player is not the current player")
	}

	enemyIndex := 0
	if gameService.Game.Players[playerIndex].Color == gameService.Game.Players[0].Color {
		enemyIndex = 1
	}

	if gameService.Game.Board.Fields[fieldIndex].Color != gameService.Game.Players[enemyIndex].Color {
		return errors.New("field is not owned by enemy")
	}

	enemyColor := gameService.Game.Players[enemyIndex].Color
	if gameService.BoardService.EnemyHasStoneOutsideMill(gameService.Game.Board, enemyColor) {
		if gameService.BoardService.IsFieldPartOfClosedMill(gameService.Game.Board, fieldIndex, enemyColor) {
			return errors.New("cannot remove stone from mill")
		}
	}

	gameService.Game.Board.Fields[fieldIndex].Color = entities.ColorUnknown

	gameService.Game.Players[enemyIndex].Stones--

	if gameService.Game.Players[enemyIndex].Stones <= 2 {
		if gameService.Game.Players[enemyIndex].Color == entities.ColorWhite {
			gameService.Game.State = entities.GameStateWinBlack
		} else {
			gameService.Game.State = entities.GameStateWinWhite
		}

		return nil
	}

	gameService.Game.State = gameService.Game.StateBeforeRemovingStone

	if gameService.Game.State == entities.GameStatePuttingStone {
		gameService.Game.Players[playerIndex].Stones++
		gameService.Game.Players[playerIndex].PuttedStones++

		// both players put 9 stones
		if gameService.Game.Players[0].PuttedStones == 9 && gameService.Game.Players[1].PuttedStones == 9 {
			gameService.Game.State = entities.GameStateMovingStone
		}
	}

	if gameService.Game.Players[enemyIndex].Stones == 3 {
		gameService.Game.Players[enemyIndex].IsJumping = true
	}

	gameService.canCurrentPlayerMove()
	gameService.nextPlayer()

	return nil
}

func (gameService *GameService) canCurrentPlayerMove() {
	player := gameService.Game.Players[gameService.Game.CurrentPlayerIndex]

	if player.IsJumping {
		return
	}

	for _, field := range gameService.Game.Board.Fields {
		if field.Color == player.Color {
			for _, neighbor := range field.Neighbors {
				if neighbor.Color == entities.ColorUnknown {
					return
				}
			}
		}
	}

	if player.Color == entities.ColorWhite {
		gameService.Game.State = entities.GameStateWinBlack
	}

	gameService.Game.State = entities.GameStateWinWhite
}

func (gameService *GameService) getPlayerIndexBySecretCode(secretCode string) int {
	for i, player := range gameService.Game.Players {
		if player.Secret == secretCode {
			return i
		}
	}
	return -1
}

func genSecretCode() string {
	// return random uuid
	id := uuid.New()

	return id.String()
}

func (gameService *GameService) GetBoard() entities.Board {
	return gameService.Game.Board
}

func NewGameService() *GameService {
	return &GameService{
		BoardService: NewBoardService(),
		Game:         entities.Game{},
	}
}
