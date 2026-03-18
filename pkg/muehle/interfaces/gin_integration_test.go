package interfaces

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRoughGameFlow_PlaceStonesAlternating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := NewClient()
	r := gin.New()
	r.Use(client.CORS)
	client.generateRouting(r)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games", nil))
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /games: %d %s", w.Code, w.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	gameID := created.ID

	postForm := func(path string, data url.Values) *httptest.ResponseRecorder {
		ww := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(ww, req)
		return ww
	}
	get := func(path string) *httptest.ResponseRecorder {
		ww := httptest.NewRecorder()
		r.ServeHTTP(ww, httptest.NewRequest(http.MethodGet, path, nil))
		return ww
	}

	base := "/games/" + gameID

	w1 := postForm(base+"/players", url.Values{"playerName": {"Weiß"}})
	if w1.Code != http.StatusOK {
		t.Fatalf("player1: %d %s", w1.Code, w1.Body.String())
	}
	var p1 struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w1.Body.Bytes(), &p1)

	w2 := postForm(base+"/players", url.Values{"playerName": {"Schwarz"}})
	if w2.Code != http.StatusOK {
		t.Fatalf("player2: %d %s", w2.Code, w2.Body.String())
	}
	var p2 struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w2.Body.Bytes(), &p2)

	// Weiß beginnt, dann abwechselnd: Felder 0..5
	for i := 0; i < 6; i++ {
		sec := p1.Secret
		if i%2 == 1 {
			sec = p2.Secret
		}
		wm := postForm(base+"/moves", url.Values{
			"action":     {"place"},
			"fieldIndex": {strconv.Itoa(i)},
			"secretCode": {sec},
		})
		if wm.Code != http.StatusOK {
			t.Fatalf("place Feld %d (secret …%s): %d %s", i, sec[len(sec)-4:], wm.Code, wm.Body.String())
		}
	}

	st := get(base + "/state")
	var state struct {
		State string `json:"state"`
	}
	json.Unmarshal(st.Body.Bytes(), &state)

	bd := get(base + "/board")
	var board struct {
		Board struct {
			Fields []struct {
				Index int    `json:"Index"`
				Color int    `json:"Color"` // JSON numbers for enum
			} `json:"Fields"`
		} `json:"board"`
	}
	json.Unmarshal(bd.Body.Bytes(), &board)
	occupied := 0
	for _, f := range board.Board.Fields {
		if f.Color != 0 { // ColorUnknown = 0
			occupied++
		}
	}
	if occupied != 6 {
		t.Errorf("erwartet 6 gesetzte Steine, got %d", occupied)
	}

	t.Logf("OK: game=%s…, State=%s, %d Steine auf dem Brett", gameID[:8], state.State, occupied)
}
