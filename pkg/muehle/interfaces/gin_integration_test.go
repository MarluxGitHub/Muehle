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
	"github.com/google/uuid"
)

func testConcurrentGamesRouter(t *testing.T) (*gin.Engine, func(string, url.Values) *httptest.ResponseRecorder, func(string) *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	client := NewClient()
	r := gin.New()
	r.Use(client.CORS)
	client.generateRouting(r)
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
	return r, postForm, get
}

func createGameID(t *testing.T, r *gin.Engine) string {
	t.Helper()
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
	return created.ID
}

func TestConcurrentGames_TwoGamesCreatedBothAddressable(t *testing.T) {
	r, _, get := testConcurrentGamesRouter(t)
	id1 := createGameID(t, r)
	id2 := createGameID(t, r)
	if id1 == id2 {
		t.Fatalf("expected two different game ids, got duplicate %s", id1)
	}
	for _, id := range []string{id1, id2} {
		w := get("/games/" + id + "/state")
		if w.Code != http.StatusOK {
			t.Fatalf("GET /games/%s/state: %d %s", id, w.Code, w.Body.String())
		}
		var st struct {
			State string `json:"state"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil {
			t.Fatal(err)
		}
		if st.State != "WaitingForPlayers" {
			t.Errorf("game %s: state %q, want WaitingForPlayers", id[:8], st.State)
		}
	}
}

func TestConcurrentGames_FiveGamesMoveInFirstOnlyOthersUnchanged(t *testing.T) {
	r, postForm, get := testConcurrentGamesRouter(t)
	ids := make([]string, 5)
	var whiteSecret string
	for i := range ids {
		ids[i] = createGameID(t, r)
		base := "/games/" + ids[i]
		w1 := postForm(base+"/players", url.Values{"playerName": {"Weiß"}})
		if w1.Code != http.StatusOK {
			t.Fatalf("players game %d: %d %s", i, w1.Code, w1.Body.String())
		}
		var p1 struct {
			Secret string `json:"secret"`
		}
		if err := json.Unmarshal(w1.Body.Bytes(), &p1); err != nil {
			t.Fatal(err)
		}
		w2 := postForm(base+"/players", url.Values{"playerName": {"Schwarz"}})
		if w2.Code != http.StatusOK {
			t.Fatalf("player2 game %d: %d %s", i, w2.Code, w2.Body.String())
		}
		if i == 0 {
			whiteSecret = p1.Secret
		}
	}
	snapshots := make([]string, 4)
	for i := 1; i < 5; i++ {
		w := get("/games/" + ids[i] + "/board")
		if w.Code != http.StatusOK {
			t.Fatalf("GET board game %d: %d", i, w.Code)
		}
		snapshots[i-1] = w.Body.String()
	}
	wm := postForm("/games/"+ids[0]+"/moves", url.Values{
		"action":     {"place"},
		"fieldIndex": {"0"},
		"secretCode": {whiteSecret},
	})
	if wm.Code != http.StatusOK {
		t.Fatalf("place in game 0: %d %s", wm.Code, wm.Body.String())
	}
	for i := 1; i < 5; i++ {
		w := get("/games/" + ids[i] + "/board")
		if w.Body.String() != snapshots[i-1] {
			t.Errorf("board of game %d changed after move in game 0", i+1)
		}
	}
}

func TestConcurrentGames_UnknownGameIdMoveReturns404AndDoesNotTouchExisting(t *testing.T) {
	r, postForm, get := testConcurrentGamesRouter(t)
	id := createGameID(t, r)
	base := "/games/" + id
	if postForm(base+"/players", url.Values{"playerName": {"A"}}).Code != http.StatusOK {
		t.Fatal("add player A")
	}
	if postForm(base+"/players", url.Values{"playerName": {"B"}}).Code != http.StatusOK {
		t.Fatal("add player B")
	}
	wb := get(base + "/board")
	if wb.Code != http.StatusOK {
		t.Fatal(wb.Code)
	}
	snap := wb.Body.String()

	unknown := uuid.New().String()
	for unknown == id {
		unknown = uuid.New().String()
	}
	wbad := postForm("/games/"+unknown+"/moves", url.Values{
		"action":     {"place"},
		"fieldIndex": {"0"},
		"secretCode": {"irrelevant"},
	})
	if wbad.Code != http.StatusNotFound {
		t.Fatalf("unknown game move: status %d want 404, body %s", wbad.Code, wbad.Body.String())
	}

	wb2 := get(base + "/board")
	if wb2.Body.String() != snap {
		t.Error("board of existing game changed after 404 on unknown game")
	}
}

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
				Index int `json:"Index"`
				Color int `json:"Color"` // JSON numbers for enum
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

func TestOpenAPIAndSwaggerRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := NewClient()
	r := gin.New()
	r.Use(c.CORS)
	c.generateRouting(r)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET /openapi.yaml: %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/yaml" {
		t.Errorf("Content-Type: got %q want application/yaml", ct)
	}
	if len(w.Body.Bytes()) < 100 || !strings.Contains(w.Body.String(), "openapi:") {
		t.Fatal("openapi.yaml body invalid")
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/swagger", nil))
	if w2.Code != http.StatusOK {
		t.Fatalf("GET /swagger: %d", w2.Code)
	}
	if !strings.Contains(w2.Body.String(), "swagger-ui") {
		t.Fatal("swagger index should reference swagger-ui")
	}

	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, httptest.NewRequest(http.MethodGet, "/swagger", nil))
	if w3.Code != http.StatusOK {
		t.Fatalf("GET /swagger: %d", w3.Code)
	}
	if !strings.Contains(w3.Body.String(), "swagger-ui") {
		t.Fatal("GET /swagger should serve Swagger UI HTML")
	}

	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, httptest.NewRequest(http.MethodGet, "/health", nil))
	if w4.Code != http.StatusOK || w4.Body.String() != "ok" {
		t.Fatalf("GET /health: %d %q", w4.Code, w4.Body.String())
	}
}
