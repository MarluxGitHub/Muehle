package entities

type Player struct {
	ID           int
	Name         string
	Secret       string
	Color        Color
	Stones       int
	PuttedStones int
	IsJumping    bool
}
