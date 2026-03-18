package services

import (
	"marluxgithub/muehle/pkg/muehle/domain/entities"
	"testing"
)

// TestFullGame_PlacingPhaseToMovingStone: leere Partie, 18 Steine setzen (inkl. Mühlen schlagen).
func TestFullGame_PlacingPhaseToMovingStone(t *testing.T) {
	gs := NewGameService()
	gs.CreateGame()
	_, err := gs.AddPlayer(entities.Player{Name: "A"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = gs.AddPlayer(entities.Player{Name: "B"})
	if err != nil {
		t.Fatal(err)
	}
	for step := 0; step < 500; step++ {
		if gs.Game.State == entities.GameStateMovingStone {
			t.Logf("Setzphase nach %d API-Schritten beendet", step)
			return
		}
		if gs.Game.State == entities.GameStateWinWhite || gs.Game.State == entities.GameStateWinBlack {
			t.Fatalf("unerwarteter Sieg in Setzphase: %s", gs.Game.State.String())
		}
		idx := gs.Game.CurrentPlayerIndex
		sec := gs.Game.Players[idx].Secret
		if gs.Game.State == entities.GameStateRemovingStone {
			removed := false
			for f := 0; f < 24; f++ {
				if gs.RemoveStone(f, sec) == nil {
					removed = true
					break
				}
			}
			if !removed {
				t.Fatalf("Schritt %d: kein gültiger Schlag (RemovingStone)", step)
			}
			continue
		}
		if gs.Game.State != entities.GameStatePuttingStone {
			t.Fatalf("Schritt %d: unerwarteter State %s", step, gs.Game.State.String())
		}
		placed := false
		for f := 0; f < 24; f++ {
			if gs.MovePutStone(f, sec) == nil {
				placed = true
				break
			}
		}
		if !placed {
			t.Fatalf("Schritt %d: kein leeres Feld für %s", step, gs.GetCurrentPlayerColor().String())
		}
	}
	t.Fatal("Setzphase endet nicht (MovingStone)")
}

// TestFullGame_CreateFakeGameWhiteWinsByTakingBlackToTwo: nahezu Endstellung, ein Schlag → Sieg Weiß.
func TestFullGame_CreateFakeGameWhiteWinsByTakingBlackToTwo(t *testing.T) {
	for from := 0; from < 24; from++ {
		for to := 0; to < 24; to++ {
			if from == to {
				continue
			}
			gs := NewGameService()
			gs.CreateFakeGame()
			if gs.MoveStone(from, to, "1") != nil {
				continue
			}
			if gs.Game.State != entities.GameStateRemovingStone {
				continue
			}
			for rm := 0; rm < 24; rm++ {
				g2 := NewGameService()
				g2.CreateFakeGame()
				if g2.MoveStone(from, to, "1") != nil {
					t.Fatal("replay move")
				}
				if g2.RemoveStone(rm, "1") != nil {
					continue
				}
				if g2.Game.State == entities.GameStateWinWhite {
					t.Logf("Sieg Weiß: Zug %d→%d, Schlag Feld %d", from, to, rm)
					return
				}
			}
		}
	}
	t.Fatal("kein Weiß-Sieg aus CreateFakeGame mit einem Zug+Schlag gefunden")
}
