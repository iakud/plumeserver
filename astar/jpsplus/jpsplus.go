package jpsplus

import (
	"image"
)

type direction int

const (
	kUpLeft direction = iota
	kUp
	kUpRight
	kLeft
	kRight
	kDownRight
	kDown
	kDownLeft
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
		kUpLeft: dirUpLeft, kUp: dirUp, kUpRight: dirUpRight, // UpLeft, Up, UpRight
		kLeft: dirLeft, kRight: dirRight, // Left, Right
		kDownLeft: dirDownLeft, kDown: dirDown, kDownRight: dirDownRight, // DownLeft, Down, DownRight
	}

	// straight & diagonal
	straightDirectionIndexes = []direction{kUp, kRight, kDown, kLeft}
	diagonalDirectionIndexes = []direction{kUpLeft, kUpRight, kDownRight, kDownLeft}

	diagonalToStraightDirection = map[direction][]direction{
		kUpLeft: {kUp, kLeft},
		kUpRight: {kUp, kRight},
		kDownRight :{kDown, kRight},
		kDownLeft: {kDown, kLeft},
	}

	jumpDirections = [8]uint8{
		kUpLeft: jumpUpLeft, kUp: jumpUp, kUpRight: jumpUpRight, // UpLeft, Up, UpRight
		kLeft: jumpLeft, kRight: jumpRight, // Left, Right
		kDownLeft: jumpDownLeft, kDown: jumpDown, kDownRight: jumpDownRight, // DownLeft, Down, DownRight
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
	directionMap = map[image.Point][]direction{
		dirUpLeft:    {kLeft, kUpLeft, kUp},
		dirUp:        {kLeft, kUpLeft, kUp, kUpRight, kRight},
		dirUpRight:   {kUp, kUpRight, kRight},
		dirRight:     {kUp, kUpRight, kRight, kDownRight, kDown},
		dirDownRight: {kRight, kDownRight, kDown},
		dirDown:      {kRight, kDownRight, kDown, kDownLeft, kLeft},
		dirDownLeft:  {kDown, kDownLeft, kLeft},
		dirLeft:      {kDown, kDownLeft, kLeft, kUpLeft, kUp},

		image.ZP: {kUpLeft, kUp, kUpRight, kRight, kDownRight, kDown, kDownLeft, kLeft},
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
			p := image.Pt(x, y)
			if !g.IsFreeAt(p) {
				continue
			}
			for _, directionIndex := range straightDirectionIndexes {
				nodes[y][x].jumpDistance[directionIndex] = straightDistance(g, jg, p, directionIndex)
			}
		}
	}
	// diagonal
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			p := image.Pt(x, y)
			if !g.IsFreeAt(p) {
				continue
			}
			for _, directionIndex := range diagonalDirectionIndexes {
				nodes[p.Y][p.X].jumpDistance[directionIndex] = diagonalDistance(g, jg, p, directionIndex, nodes)
			}
		}
	}
	return &graph{nodes: nodes}
}

func straightDistance(g Grid, jg jumpGrid, p image.Point, directionIndex direction) int {
	d := directions[directionIndex]
	jump := jumpDirections[directionIndex]
	var distance int = 0
	for n := p.Add(d); g.IsFreeAt(n); n = n.Add(d) {
		distance++
		if jg[n.Y][n.X]&jump != 0 {
			return distance
		}
	}
	return -distance
}

func diagonalDistance(g Grid, jg jumpGrid, p image.Point, directionIndex direction, nodes [][]*node) int {
	d := directions[directionIndex]
	jump := jumpDirections[directionIndex]

	// is near by block
	if !g.IsFreeAt(p.Add(image.Pt(d.X, 0))) || !g.IsFreeAt(p.Add(image.Pt(0, d.Y))) {
		return 0
	}
	var distance int = 0
	for n := p.Add(d); g.IsFreeAt(n); n = n.Add(d) {
		distance++
		if jg[n.Y][n.X]&jump != 0 {
			return distance
		}
		for _, straightDirectionIndex := range diagonalToStraightDirection[directionIndex] {
			if nodes[n.Y][n.X].jumpDistance[straightDirectionIndex] > 0 {
				return distance
			}
		}
	}
	return -distance
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
