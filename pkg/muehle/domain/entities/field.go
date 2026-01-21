package entities

type Field struct {
	Index     int
	Color     Color
	Neighbors []Field `json:"-"`
}
