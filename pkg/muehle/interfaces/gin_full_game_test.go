package interfaces

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

// TestE2E_HTTP_FullPlacingPhaseThenAtLeastOneMove: kompletter Ablauf über HTTP wie ein Client.
func TestE2E_HTTP_FullPlacingPhaseThenAtLeastOneMove(t *testing.T) {
	r, postForm, get := testConcurrentGamesRouter(t)
	id := createGameID(t, r)
	base := "/games/" + id

	w1 := postForm(base+"/players", url.Values{"playerName": {"Weiß"}})
	if w1.Code != http.StatusOK {
		t.Fatal(w1.Body.String())
	}
	var pW struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w1.Body.Bytes(), &pW)
	w2 := postForm(base+"/players", url.Values{"playerName": {"Schwarz"}})
	if w2.Code != http.StatusOK {
		t.Fatal(w2.Body.String())
	}
	var pB struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w2.Body.Bytes(), &pB)

	gameState := func() string {
		w := get(base + "/state")
		var v struct {
			State string `json:"state"`
		}
		json.Unmarshal(w.Body.Bytes(), &v)
		return v.State
	}
	currentColor := func() string {
		w := get(base + "/current-player")
		var v struct {
			Color string `json:"color"`
		}
		json.Unmarshal(w.Body.Bytes(), &v)
		return v.Color
	}
	secret := func() string {
		if currentColor() == "Black" {
			return pB.Secret
		}
		return pW.Secret
	}

	for step := 0; step < 500; step++ {
		st := gameState()
		if st == "MovingStone" {
			break
		}
		if st == "WinWhite" || st == "WinBlack" {
			t.Fatalf("früher Sieg in Setzphase: %s", st)
		}
		sec := secret()
		if st == "RemovingStone" {
			ok := false
			for f := 0; f < 24; f++ {
				w := postForm(base+"/moves", url.Values{
					"action":     {"remove"},
					"fieldIndex": {strconv.Itoa(f)},
					"secretCode": {sec},
				})
				if w.Code == http.StatusOK {
					ok = true
					break
				}
			}
			if !ok {
				t.Fatalf("Schritt %d: Schlagen unmöglich", step)
			}
			continue
		}
		if st != "PuttingStone" {
			t.Fatalf("State %s", st)
		}
		ok := false
		for f := 0; f < 24; f++ {
			w := postForm(base+"/moves", url.Values{
				"action":     {"place"},
				"fieldIndex": {strconv.Itoa(f)},
				"secretCode": {sec},
			})
			if w.Code == http.StatusOK {
				ok = true
				break
			}
		}
		if !ok {
			t.Fatalf("Schritt %d: kein Platz für %s: %s", step, currentColor(), gameState())
		}
	}
	if gameState() != "MovingStone" {
		t.Fatalf("erwartet MovingStone, got %s", gameState())
	}

	// mindestens einen Zug in Ziehphase
	sec := secret()
	moved := false
	for from := 0; from < 24; from++ {
		for to := 0; to < 24; to++ {
			w := postForm(base+"/moves", url.Values{
				"action":       {"move"},
				"fieldIndex":   {strconv.Itoa(from)},
				"toFieldIndex": {strconv.Itoa(to)},
				"secretCode":   {sec},
			})
			if w.Code == http.StatusOK {
				moved = true
				goto movedOK
			}
		}
	}
movedOK:
	if !moved {
		t.Fatal("kein legaler Zug in Ziehphase (oder sofortiger Sieg ohne Zug)")
	}
	st := gameState()
	if st != "MovingStone" && st != "RemovingStone" && st != "WinWhite" && st != "WinBlack" {
		t.Errorf("nach Zug unerwarteter State %s", st)
	}
	t.Logf("OK: Setzphase + erster Zug, State=%s", st)
}
