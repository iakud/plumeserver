package jpsplus

import (
	"image"
)

type direction int

const (
	kUpLeft direction = iota
	kUp
	kUpRight
	kRight
	kDownRight
	kDown
	kDownLeft
	kLeft
	kDirectionMax
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
	// straight direction vectors
	vectorUp    = image.Point{X: 0, Y: -1}
	vectorRight = image.Point{X: 1, Y: 0}
	vectorDown  = image.Point{X: 0, Y: 1}
	vectorLeft  = image.Point{X: -1, Y: 0}

	// diagonal direction vectors
	vectorUpLeft    = image.Point{X: -1, Y: -1}
	vectorUpRight   = image.Point{X: 1, Y: -1}
	vectorDownRight = image.Point{X: 1, Y: 1}
	vectorDownLeft  = image.Point{X: -1, Y: 1}

	// direction vectors
	directionVectors = []image.Point{
		kUpLeft:    image.Pt(-1, -1), // UpLeft
		kUp:        image.Pt(0, -1),  // Up
		kUpRight:   image.Pt(1, -1),  // UpRight
		kRight:     image.Pt(1, 0),   // Right
		kDownRight: image.Pt(1, 1),   // DownRight
		kDown:      image.Pt(0, 1),   // Down
		kDownLeft:  image.Pt(-1, 1),  // DownLeft
		kLeft:      image.Pt(-1, 0),  // Left
	}

	// straight & diagonal
	straightDirection = []direction{kUp, kRight, kDown, kLeft}
	diagonalDirection = []direction{kUpLeft, kUpRight, kDownRight, kDownLeft}

	diagonalToStraightDirection = map[direction][]direction{
		kUpLeft:    {kUp, kLeft},    // UpLeft
		kUpRight:   {kUp, kRight},   // UpRight
		kDownRight: {kDown, kRight}, // DownRight
		kDownLeft:  {kDown, kLeft},  // DownLeft
	}
	// jump directions
	jumpDirections = [8]uint8{
		kUpLeft:    jumpUpLeft,    // UpLeft
		kUp:        jumpUp,        // Up
		kUpRight:   jumpUpRight,   // UpRight
		kRight:     jumpRight,     // Right
		kDownRight: jumpDownRight, // DownRight
		kDown:      jumpDown,      // Down
		kDownLeft:  jumpDownLeft,  // DownLeft
		kLeft:      jumpLeft,      // Left
	}

	vectorDirections = map[image.Point][]direction{
		vectorUpLeft:    {kLeft, kUpLeft, kUp},
		vectorUp:        {kLeft, kUpLeft, kUp, kUpRight, kRight},
		vectorUpRight:   {kUp, kUpRight, kRight},
		vectorRight:     {kUp, kUpRight, kRight, kDownRight, kDown},
		vectorDownRight: {kRight, kDownRight, kDown},
		vectorDown:      {kRight, kDownRight, kDown, kDownLeft, kLeft},
		vectorDownLeft:  {kDown, kDownLeft, kLeft},
		vectorLeft:      {kDown, kDownLeft, kLeft, kUpLeft, kUp},

		image.ZP: {kUpLeft, kUp, kUpRight, kRight, kDownRight, kDown, kDownLeft, kLeft},
	}
)

type Grid interface {
	IsFreeAt(n image.Point) bool
}

type jumpGrid []uint8

func newJumpGrid(g Grid, size image.Point) jumpGrid {
	jg := make([]uint8, size.Y*size.X)
	for i:=0; i<len(jg); i++{
		p:=image.Pt(i%size.X, i/size.X)
		if !g.IsFreeAt(p) {
			continue
		}
		var jump uint8
		for _, d := range straightDirection {
			if !isJumpPoint(g, p, directionVectors[d]) {
				continue
			}
			jump |= jumpDirections[d]
		}
		jg[i] = jump
	}
	return jg
}

func isJumpPoint(g Grid, p image.Point, v image.Point) bool {
	if !g.IsFreeAt(p.Sub(v)) {
		return false
	}
	// forced neighbour
	if g.IsFreeAt(image.Pt(p.X+v.Y, p.Y+v.X)) && !g.IsFreeAt(image.Pt(p.X+v.Y, p.Y+v.X).Sub(v)) {
		return true
	}
	if g.IsFreeAt(image.Pt(p.X-v.Y, p.Y-v.X)) && !g.IsFreeAt(image.Pt(p.X-v.Y, p.Y-v.X).Sub(v)) {
		return true
	}
	return false
}

type node [8]int

// FIXME: Neighbors(n, p)
func (g *graph) Neighbors(n, p, goal image.Point) []image.Point {
	jumpNodes := make([]image.Point, 0)
	jumpDistance := g.nodes[n.Y*g.size.X + n.X]
	// direction vector
	for _, d := range vectorDirections[image.Pt((n.X-p.X)/max(abs(n.X-p.X), 1), (n.Y-p.Y)/max(abs(n.Y-p.Y), 1))] {
		distance := jumpDistance[d]
		if distance == 0 {
			continue
		}
		v := directionVectors[d]
		to := n.Add(v.Mul(abs(distance)))
		// 点在线上
		if (goal.Y - n.Y) * v.X == (goal.X - n.X) * v.Y &&
			(n.X < goal.X && goal.X <= to.X || to.X <= goal.X && goal.X < n.X) &&
			(n.Y < goal.Y && goal.Y <= to.Y || to.Y <= goal.Y && goal.Y < n.Y) {
			jumpNodes = append(jumpNodes, goal)
		} else {
			jumpNodes = append(jumpNodes, to)
		}
	}
	return jumpNodes
}

type graph struct {
	size image.Point
	nodes []node
}

func NewGraph(g Grid, size image.Point) *graph {
	jg := newJumpGrid(g, size)
	nodes := make([]node, size.Y*size.X)

	// straight
	for i:=0; i<len(nodes); i++ {
		p:=image.Pt(i%size.X, i/size.X)
		if !g.IsFreeAt(p) {
			continue
		}
		node := &nodes[i]
		for _, d := range straightDirection {
			node[d] = straightDistance(g, size, jg, p, d)
		}
	}

	// diagonal
	for i:=0; i<len(nodes); i++ {
		p:=image.Pt(i%size.X, i/size.X)
		if !g.IsFreeAt(p) {
			continue
		}
		node := &nodes[i]
		for _, d := range diagonalDirection {
			node[d] = diagonalDistance(g, size, jg, p, d, nodes)
		}
	}

	return &graph{size: size, nodes: nodes}
}

func straightDistance(g Grid, size image.Point, jg jumpGrid, p image.Point, d direction) int {
	v := directionVectors[d]
	jump := jumpDirections[d]
	var distance int = 0
	for n := p.Add(v); g.IsFreeAt(n); n = n.Add(v) {
		distance++
		if jg[n.Y*size.X + n.X]&jump != 0 {
			return distance
		}
	}
	return -distance
}

func diagonalDistance(g Grid, size image.Point, jg jumpGrid, p image.Point, d direction, nodes []node) int {
	v := directionVectors[d]
	jump := jumpDirections[d]

	// is near by obstacle
	if !g.IsFreeAt(p.Add(image.Pt(v.X, 0))) || !g.IsFreeAt(p.Add(image.Pt(0, v.Y))) {
		return 0
	}
	var distance int = 0
	for n := p.Add(v); g.IsFreeAt(n); n = n.Add(v) {
		distance++
		if jg[n.Y*size.X + n.X]&jump != 0 {
			return distance
		}
		node := &nodes[n.Y*size.X + n.X]
		for _, straightDirection := range diagonalToStraightDirection[d] {
			if node[straightDirection] > 0 {
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
