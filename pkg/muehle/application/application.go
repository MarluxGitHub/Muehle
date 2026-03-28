package application

import (
	"marluxgithub/muehle/pkg/muehle/domain/entities"
	"marluxgithub/muehle/pkg/muehle/domain/services"
)

type Application struct {
	GameService *services.GameService
}

func (application *Application) CreateGame() error {
	application.GameService.CreateFakeGame()
	return nil
}

func (application *Application) AddPlayer(playerName string) (secret string, color string, err error) {
	player := entities.Player{
		Name: playerName,
	}
	s, c, e := application.GameService.AddPlayer(player)
	if e != nil {
		return "", "", e
	}
	return s, c.String(), nil
}

func (application *Application) MovePutStone(fieldIndex int, secretCode string) error {
	return application.GameService.MovePutStone(fieldIndex, secretCode)
}

func (application *Application) MoveStone(fromFieldIndex, toFieldIndex int, secretCode string) error {
	return application.GameService.MoveStone(fromFieldIndex, toFieldIndex, secretCode)
}

func (application *Application) RemoveStone(fieldIndex int, secretCode string) error {
	return application.GameService.RemoveStone(fieldIndex, secretCode)
}

func (application *Application) GetGameState() string {
	return application.GameService.GetGameState().String()
}

func (application *Application) GetCurrentPlayerColor() string {
	return application.GameService.GetCurrentPlayerColor().String()
}

func (application *Application) GetBoard() entities.Board {
	return application.GameService.GetBoard()
}

func NewApplication() *Application {
	return &Application{
		GameService: services.NewGameService(),
	}
}
