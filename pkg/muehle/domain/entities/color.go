package entities

type Color int

const (
	ColorUnknown Color = iota
	ColorWhite
	ColorBlack
)

func (color Color) String() string {
	return []string{
		"Unknown",
		"White",
		"Black",
	}[color]
}
