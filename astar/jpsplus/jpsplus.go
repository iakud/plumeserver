package jpsplus

import (
	"image"
)

type directionIndex int

const (
	indexUpLeft directionIndex = iota
	indexUp
	indexUpRight
	indexLeft
	indexRight
	indexDownRight
	indexDown
	indexDownLeft
)

const (
	// straight jump direction
	jumpUp = 1 << iota
	jumpRight
	jumpDown
	jumpLeft
	// diagonal jump direction
	jumpUpLeft    = jumpUp | jumpLeft
	jumpUpRight   = jumpUp | jumpRight
	jumpDownRight = jumpDown | jumpRight
	jumpDownLeft  = jumpDown | jumpLeft
)

var (
	// straight direction
	dirUp    = image.Point{X: 0, Y: -1}
	dirRight = image.Point{X: 1, Y: 0}
	dirDown  = image.Point{X: 0, Y: 1}
	dirLeft  = image.Point{X: -1, Y: 0}

	// diagonal direction
	dirUpLeft    = image.Point{X: -1, Y: -1}
	dirUpRight   = image.Point{X: 1, Y: -1}
	dirDownRight = image.Point{X: 1, Y: 1}
	dirDownLeft  = image.Point{X: -1, Y: 1}

	// directions
	directions = []image.Point{
		indexUpLeft: dirUpLeft, indexUp: dirUp, indexUpRight: dirUpRight, // UpLeft, Up, UpRight
		indexLeft: dirLeft, indexRight: dirRight, // Left, Right
		indexDownLeft: dirDownLeft, indexDown: dirDown, indexDownRight: dirDownRight, // DownLeft, Down, DownRight
	}

	// straight & diagonal
	straightDirectionIndexes = []directionIndex{indexUp, indexRight, indexDown, indexLeft}
	diagonalDirectionIndexes = []directionIndex{indexUpLeft, indexUpRight, indexDownRight, indexDownLeft}

	diagonalToStraightDirectionIndexes = [][]directionIndex{
		{indexUp, indexLeft},
		{indexUp, indexRight},
		{indexDown, indexRight},
		{indexDown, indexLeft},
	}

	jumpDirections = [8]uint8{
		indexUpLeft: jumpUpLeft, indexUp: jumpUp, indexUpRight: jumpUpRight,
		indexLeft: jumpLeft, indexRight: jumpRight,
		indexDownLeft: jumpDownLeft, indexDown: jumpDown, indexDownRight: jumpDownRight,
	}
)

type node struct {
	jumpDistance [8]int
}

type Grid interface {
	Size() image.Point
	IsFreeAt(n image.Point) bool
}

type jumpGrid [][]uint8

func newJumpGrid(g Grid) jumpGrid {
	size := g.Size()
	jg := make([][]uint8, size.Y)
	for y := 0; y < size.Y; y++ {
		jg[y] = make([]uint8, size.X)
		for x := 0; x < size.X; x++ {
			p := image.Pt(x, y)
			if !g.IsFreeAt(p) {
				continue
			}
			var jump uint8
			for _, directionIndex := range straightDirectionIndexes {
				d := directions[directionIndex]
				if !isJumpPoint(g, p, d) {
					continue
				}
				jump |= jumpDirections[directionIndex]
			}
			jg[y][x] = jump
		}
	}
	return jg
}

func isJumpPoint(g Grid, p image.Point, d image.Point) bool {
	if !g.IsFreeAt(p.Sub(d)) {
		return false
	}
	// forced neighbour
	if g.IsFreeAt(image.Pt(p.X+d.Y, p.Y+d.X)) && !g.IsFreeAt(image.Pt(p.X+d.Y, p.Y+d.X).Sub(d)) {
		return true
	}
	if g.IsFreeAt(image.Pt(p.X-d.Y, p.Y-d.X)) && !g.IsFreeAt(image.Pt(p.X-d.Y, p.Y-d.X).Sub(d)) {
		return true
	}
	return false
}

var (
	directionMap = map[image.Point][]directionIndex{
		dirUpLeft:    {indexLeft, indexUpLeft, indexUp},
		dirUp:        {indexLeft, indexUpLeft, indexUp, indexUpRight, indexRight},
		dirUpRight:   {indexUp, indexUpRight, indexRight},
		dirRight:     {indexUp, indexUpRight, indexRight, indexDownRight, indexDown},
		dirDownRight: {indexRight, indexDownRight, indexDown},
		dirDown:      {indexRight, indexDownRight, indexDown, indexDownLeft, indexLeft},
		dirDownLeft:  {indexDown, indexDownLeft, indexLeft},
		dirLeft:      {indexDown, indexDownLeft, indexLeft, indexUpLeft, indexUp},

		image.ZP: {indexUpLeft, indexUp, indexUpRight, indexRight, indexDownRight, indexDown, indexDownLeft, indexLeft},
	}
)

type graph struct {
	nodes [][]*node
}

func (g *graph) Neighbors(n, from, goal image.Point) []image.Point {
	jumpNodes := make([]image.Point, 0)
	jumpDistance := g.nodes[n.Y][n.X].jumpDistance
	dir := image.Pt((n.X-from.X)/max(abs(n.X-from.X), 1), (n.Y-from.Y)/max(abs(n.Y-from.Y), 1))
	for _, directionIndex := range directionMap[dir] {
		distance := jumpDistance[directionIndex]
		if distance == 0 {
			continue
		}
		dir := directions[directionIndex]
		if distance < 0 {
			distance = -distance
		}
		to := n.Add(dir.Mul(distance))
		endDir := goal.Sub(n)
		// 方向平行
		if endDir.Y*dir.X == endDir.X*dir.Y &&
			(n.X < goal.X && goal.X <= to.X || to.X <= goal.X && goal.X < n.X) &&
			(n.Y < goal.Y && goal.Y <= to.Y || to.Y <= goal.Y && goal.Y < n.Y) {
			// 点在线上
			jumpNodes = append(jumpNodes, goal)
		} else {
			jumpNodes = append(jumpNodes, to)
		}
	}
	return jumpNodes
}

func NewGraph(g Grid) *graph {
	size := g.Size()
	jg := newJumpGrid(g)
	nodes := make([][]*node, size.Y)
	for y := 0; y < size.Y; y++ {
		nodes[y] = make([]*node, size.X)
		for x := 0; x < size.X; x++ {
			nodes[y][x] = &node{}
		}
	}
	// straight
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			nodes[y][x].jumpDistance = straightDirectionDistance(g, jg, image.Pt(x, y))
		}
	}
	// diagonal
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			diagonalDirectionDistance(g, jg, nodes, image.Pt(x, y))
		}
	}
	return &graph{nodes: nodes}
}

func straightDirectionDistance(g Grid, jg jumpGrid, p image.Point) [8]int {
	jumpDistance := [8]int{}
	if !g.IsFreeAt(p) {
		return jumpDistance
	}
	for _, directionIndex := range straightDirectionIndexes {
		dir := directions[directionIndex]
		var distance int = 0
		for next := p.Add(dir); g.IsFreeAt(next); next = next.Add(dir) {
			distance--
			if jg[next.Y][next.X]&jumpDirections[directionIndex] != 0 {
				distance = -distance
				break
			}
		}
		jumpDistance[directionIndex] = distance
	}
	return jumpDistance
}

func diagonalDirectionDistance(g Grid, jg jumpGrid, nodes [][]*node, p image.Point) {
	if !g.IsFreeAt(p) {
		return
	}
	for i, directionIndex := range diagonalDirectionIndexes {
		dir := directions[directionIndex]
		// is near by block
		if !g.IsFreeAt(p.Add(image.Pt(dir.X, 0))) || !g.IsFreeAt(p.Add(image.Pt(0, dir.Y))) {
			continue
		}
		var distance int = 0
		for next := p.Add(dir); g.IsFreeAt(next); next = next.Add(dir) {
			distance--
			if jg[next.Y][next.X]&jumpDirections[directionIndex] != 0 {
				distance = -distance
				break
			}
			var found bool
			for _, straightDirectionIndex := range diagonalToStraightDirectionIndexes[i] {
				if nodes[next.Y][next.X].jumpDistance[straightDirectionIndex] > 0 {
					found = true
					break
				}
			}
			if found {
				distance = -distance
				break
			}
		}
		nodes[p.Y][p.X].jumpDistance[directionIndex] = distance
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
