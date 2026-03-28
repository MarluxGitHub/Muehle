package services

import (
	"marluxgithub/muehle/pkg/muehle/domain/entities"
	"testing"
)

func TestRemoveStone_FromMillRejectedWhenEnemyHasFreeStone(t *testing.T) {
	gs := NewGameService()
	gs.CreateGame()
	sw, _, _ := gs.AddPlayer(entities.Player{Name: "W"})
	gs.AddPlayer(entities.Player{Name: "B"})
	for i := range gs.Game.Board.Fields {
		gs.Game.Board.Fields[i].Color = entities.ColorUnknown
	}
	for _, i := range []int{9, 10, 11, 12} {
		gs.Game.Board.Fields[i].Color = entities.ColorBlack
	}
	gs.Game.Players[1].Stones = 4
	gs.Game.Players[1].PuttedStones = 9
	gs.Game.State = entities.GameStateRemovingStone
	gs.Game.StateBeforeRemovingStone = entities.GameStateMovingStone
	gs.Game.CurrentPlayerIndex = 0

	err := gs.RemoveStone(9, sw)
	if err == nil || err.Error() != "cannot remove stone from mill" {
		t.Fatalf("expected cannot remove stone from mill, got %v", err)
	}
	if gs.Game.Board.Fields[9].Color != entities.ColorBlack {
		t.Fatal("board must stay unchanged on error")
	}

	err = gs.RemoveStone(12, sw)
	if err != nil {
		t.Fatalf("remove free stone: %v", err)
	}
	if gs.Game.Board.Fields[12].Color != entities.ColorUnknown {
		t.Fatal("field 12 should be empty")
	}
}

func TestRemoveStone_FromMillWhenAllEnemyStonesInMills(t *testing.T) {
	gs := NewGameService()
	gs.CreateGame()
	sw, _, _ := gs.AddPlayer(entities.Player{Name: "W"})
	gs.AddPlayer(entities.Player{Name: "B"})
	for i := range gs.Game.Board.Fields {
		gs.Game.Board.Fields[i].Color = entities.ColorUnknown
	}
	for _, i := range []int{9, 10, 11} {
		gs.Game.Board.Fields[i].Color = entities.ColorBlack
	}
	gs.Game.Players[1].Stones = 3
	gs.Game.Players[1].PuttedStones = 9
	gs.Game.State = entities.GameStateRemovingStone
	gs.Game.StateBeforeRemovingStone = entities.GameStateMovingStone
	gs.Game.CurrentPlayerIndex = 0

	err := gs.RemoveStone(10, sw)
	if err != nil {
		t.Fatalf("all in mill, remove allowed: %v", err)
	}
	if gs.Game.Board.Fields[10].Color != entities.ColorUnknown {
		t.Fatal("stone 10 removed")
	}
}

func TestRemoveStone_NoWinWhileEnemyStillHasStonesToPlace(t *testing.T) {
	gs := NewGameService()
	gs.CreateGame()
	sw, _, _ := gs.AddPlayer(entities.Player{Name: "W"})
	_, _, _ = gs.AddPlayer(entities.Player{Name: "B"})
	for i := range gs.Game.Board.Fields {
		gs.Game.Board.Fields[i].Color = entities.ColorUnknown
	}
	gs.Game.Board.Fields[20].Color = entities.ColorBlack
	gs.Game.Players[1].Stones = 3
	gs.Game.Players[1].PuttedStones = 3
	gs.Game.CurrentPlayerIndex = 0
	gs.Game.State = entities.GameStateRemovingStone
	gs.Game.StateBeforeRemovingStone = entities.GameStatePuttingStone

	if err := gs.RemoveStone(20, sw); err != nil {
		t.Fatal(err)
	}
	if gs.Game.State == entities.GameStateWinWhite || gs.Game.State == entities.GameStateWinBlack {
		t.Fatalf("no win while opponent still has placements, got %s", gs.Game.State.String())
	}
}

func TestCanCurrentPlayerMove_PuttingStoneDoesNotAwardWin(t *testing.T) {
	gs := NewGameService()
	gs.CreateGame()
	_, _, _ = gs.AddPlayer(entities.Player{Name: "W"})
	_, _, _ = gs.AddPlayer(entities.Player{Name: "B"})
	for i := range gs.Game.Board.Fields {
		gs.Game.Board.Fields[i].Color = entities.ColorUnknown
	}
	gs.Game.Board.Fields[0].Color = entities.ColorWhite
	gs.Game.Board.Fields[1].Color = entities.ColorBlack
	gs.Game.Board.Fields[9].Color = entities.ColorBlack
	gs.Game.CurrentPlayerIndex = 0
	gs.Game.State = entities.GameStatePuttingStone

	gs.canCurrentPlayerMove()
	if gs.Game.State == entities.GameStateWinBlack || gs.Game.State == entities.GameStateWinWhite {
		t.Fatalf("PuttingStone must not use move-immobilization win, got %s", gs.Game.State.String())
	}
}
