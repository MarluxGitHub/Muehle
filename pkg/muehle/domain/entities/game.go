package entities

type Game struct {
	ID      int
	Board   Board
	Players []Player

	CurrentPlayerIndex       int
	State                    GameState
	StateBeforeRemovingStone GameState
}

type GameState int

const (
	GameStateUnknown GameState = iota
	GameStateWaitingForPlayers

	GameStatePuttingStone
	GameStateMovingStone

	GameStateRemovingStone

	GameStateWinWhite
	GameStateWinBlack
)

func (gameState GameState) String() string {
	return []string{
		"Unknown",
		"WaitingForPlayers",
		"PuttingStone",
		"MovingStone",
		"RemovingStone",
		"WinWhite",
		"WinBlack",
	}[gameState]
}
