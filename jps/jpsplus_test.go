package astar

import (
	"image"
)

type DirectionIdx int

const (
	IdxUpLeft DirectionIdx = iota
	IdxUp
	IdxUpRight
	IdxRight
	IdxDownRight
	IdxDown
	IdxDownLeft
	IdxLeft
	IdxAll
)

var (
	directionMap = map[DirectionIdx][]DirectionIdx{
		IdxDown:      {IdxLeft, IdxDownLeft, IdxDown, IdxDownRight, IdxRight},
		IdxLeft:      {IdxLeft, IdxDownLeft, IdxDown, IdxUpLeft, IdxUp},
		IdxRight:     {IdxRight, IdxDownRight, IdxDown, IdxUpRight, IdxUp},
		IdxUp:        {IdxLeft, IdxUpLeft, IdxUp, IdxUpRight, IdxRight},
		IdxUpRight:   {IdxUpRight, IdxUp, IdxRight},
		IdxDownRight: {IdxRight, IdxDownRight, IdxDown},
		IdxDownLeft:  {IdxLeft, IdxDownLeft, IdxDown},
		IdxUpLeft:    {IdxLeft, IdxUpLeft, IdxUp},
		IdxAll:       {IdxUp, IdxDown, IdxLeft, IdxRight, IdxUpLeft, IdxUpRight, IdxDownLeft, IdxDownRight},
	}

	straightDir = [4]DirectionIdx{IdxUp, IdxRight, IdxDown, IdxLeft}
	leanDir     = [4]DirectionIdx{IdxUpLeft, IdxUpRight, IdxDownRight, IdxDownLeft}

	direction8 = [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}}

	direction = [8]image.Point{
		image.Point{X: -1, Y: -1}, // UpLeft
		image.Point{X: 0, Y: -1},  // Up
		image.Point{X: 1, Y: -1},  // UpRight
		image.Point{X: 1, Y: 0},   // Right
		image.Point{X: 1, Y: 1},   // DownRight
		image.Point{X: 0, Y: 1},   // Down
		image.Point{X: -1, Y: 1},  // DownLeft
		image.Point{X: -1, Y: 0},  // Left
	}
)

type jpsPlusGraph [][]*NodePlus

type jpsGraph struct {
}

func (f jpsGraph) Neighbours(p image.Point) []image.Point {
	for _, direction := range directionMap[curr.direction] {

	}
	return nil
}
