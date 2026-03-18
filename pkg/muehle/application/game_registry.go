package application

import (
	"sync"

	"github.com/google/uuid"
)

// GameRegistry verwaltet mehrere Spiel-Sessions (In-Memory), je UUID eine Application.
type GameRegistry struct {
	mu    sync.RWMutex
	games map[string]*Application
}

func NewGameRegistry() *GameRegistry {
	return &GameRegistry{
		games: make(map[string]*Application),
	}
}

// CreateGame legt eine neue leere Partie an (Warten auf Spieler)—REST-Ablauf wie Spec/quickstart.
func (r *GameRegistry) CreateGame() (string, error) {
	app := NewApplication()
	app.GameService.CreateGame()
	id := uuid.New().String()
	r.mu.Lock()
	r.games[id] = app
	r.mu.Unlock()
	return id, nil
}

// Get liefert die Application für eine gültige UUID, die in der Registry existiert.
// Ungültiges UUID-Format oder unbekannte ID → ok == false (ohne Panic).
func (r *GameRegistry) Get(id string) (*Application, bool) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, false
	}
	key := parsed.String()
	r.mu.RLock()
	defer r.mu.RUnlock()
	app, ok := r.games[key]
	return app, ok
}
