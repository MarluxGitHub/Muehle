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
