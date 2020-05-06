package common

type Direction int

const (
	East Direction = iota
	West
	North
	South
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)

type Position struct {
	X float64
	Y float64
}
