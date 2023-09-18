package jps

import (
	"image"
	"log"
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

const (
	JPSPlusUp = 1 << iota
	JPSPlusRight
	JPSPlusDown
	JPSPlusLeft
)

var (
	jumpDirections = [8]uint8{
		JPSPlusUp | JPSPlusLeft,    // UpLeft
		JPSPlusUp,                  // Up
		JPSPlusUp | JPSPlusRight,   // UpRight
		JPSPlusRight,               // Right
		JPSPlusDown | JPSPlusRight, // DownRight
		JPSPlusDown,                // Down
		JPSPlusDown | JPSPlusLeft,  // DownLeft
		JPSPlusLeft,                // Left
	}

	UpLeft    = image.Point{X: -1, Y: -1}
	Up        = image.Point{X: 0, Y: -1}
	UpRight   = image.Point{X: 1, Y: -1}
	Right     = image.Point{X: 1, Y: 0}
	DownRight = image.Point{X: 1, Y: 1}
	Down      = image.Point{X: 0, Y: 1}
	DownLeft  = image.Point{X: -1, Y: 1}
	Left      = image.Point{X: -1, Y: 0}

	directions = [8]image.Point{
		UpLeft,    // UpLeft
		Up,        // Up
		UpRight,   // UpRight
		Right,     // Right
		DownRight, // DownRight
		Down,      // Down
		DownLeft,  // DownLeft
		Left,      // Left
	}

	straightDir       = [4]DirectionIdx{IdxUp, IdxRight, IdxDown, IdxLeft}
	leanDir           = [4]DirectionIdx{IdxUpLeft, IdxUpRight, IdxDownRight, IdxDownLeft}
	leanToStraightDir = map[DirectionIdx][2]DirectionIdx{
		IdxUpLeft:    {IdxLeft, IdxUp},
		IdxUpRight:   {IdxUp, IdxRight},
		IdxDownRight: {IdxRight, IdxDown},
		IdxDownLeft:  {IdxDown, IdxLeft},
	}
)

type node struct {
	image.Point
	originVal    uint8
	jumpDistance [8]int
}

func (n *node) printDis() {
	log.Println(n.jumpDistance[IdxUpLeft], n.jumpDistance[IdxUp], n.jumpDistance[IdxUpRight])
	log.Println(n.jumpDistance[IdxLeft], " ", n.jumpDistance[IdxRight])
	log.Println(n.jumpDistance[IdxDownLeft], n.jumpDistance[IdxDown], n.jumpDistance[IdxDownRight])
}

type JPSGraph interface {
	Size() image.Point
	IsBlock(n image.Point) bool
}

type jpGraph [][]uint8

func newJumpPointGraph(g JPSGraph) jpGraph {
	size := g.Size()
	jpGraph := make([][]uint8, size.Y)
	for y := 0; y < size.Y; y++ {
		jpGraph[y] = make([]uint8, size.X)
		for x := 0; x < size.X; x++ {
			p := image.Pt(x, y)
			if g.IsBlock(p) {
				continue
			}
			var direction uint8
			for _, directionIdx := range straightDir {
				dir := directions[directionIdx]
				if !isJumpPoint(g, p, dir) {
					continue
				}
				direction |= jumpDirections[directionIdx]
			}
			jpGraph[y][x] = direction
		}
	}
	return jpGraph
}

func isJumpPoint(g JPSGraph, p image.Point, dir image.Point) bool {
	if g.IsBlock(p.Sub(dir)) {
		return false
	}
	// forced neighbour
	if !g.IsBlock(image.Pt(p.X+dir.Y, p.Y+dir.X)) && g.IsBlock(image.Pt(p.X+dir.Y, p.Y+dir.X).Sub(dir)) {
		return true
	}
	if !g.IsBlock(image.Pt(p.X-dir.Y, p.Y-dir.X)) && g.IsBlock(image.Pt(p.X-dir.Y, p.Y-dir.X).Sub(dir)) {
		return true
	}
	return false
}

type jpsPlusGraph struct {
	nodes [][]*node
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

var (
	dirMap = map[image.Point][]DirectionIdx{
		Down:      {IdxLeft, IdxDownLeft, IdxDown, IdxDownRight, IdxRight},
		Left:      {IdxLeft, IdxDownLeft, IdxDown, IdxUpLeft, IdxUp},
		Right:     {IdxRight, IdxDownRight, IdxDown, IdxUpRight, IdxUp},
		Up:        {IdxLeft, IdxUpLeft, IdxUp, IdxUpRight, IdxRight},
		UpRight:   {IdxUpRight, IdxUp, IdxRight},
		DownRight: {IdxRight, IdxDownRight, IdxDown},
		DownLeft:  {IdxLeft, IdxDownLeft, IdxDown},
		UpLeft:    {IdxLeft, IdxUpLeft, IdxUp},
		image.ZP:  {IdxUp, IdxDown, IdxLeft, IdxRight, IdxUpLeft, IdxUpRight, IdxDownLeft, IdxDownRight},
	}
)

func (g *jpsPlusGraph) Neighbors(n, from, goal image.Point) []image.Point {
	jumpNodes := make([]image.Point, 0)
	jumpDistance := g.nodes[n.Y][n.X].jumpDistance
	dir := image.Pt((n.X-from.X)/max(abs(n.X-from.X), 1), (n.Y-from.Y)/max(abs(n.Y-from.Y), 1))
	for _, direction := range dirMap[dir] {
		distance := jumpDistance[direction]
		if distance == 0 {
			continue
		}
		dir := directions[direction]
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

func newJpsPlusGraph(g JPSGraph) *jpsPlusGraph {
	size := g.Size()
	jpGraph := newJumpPointGraph(g)
	nodes := make([][]*node, size.Y)
	for y := 0; y < size.Y; y++ {
		nodes[y] = make([]*node, size.X)
		for x := 0; x < size.X; x++ {
			p := image.Pt(x, y)
			node := &node{Point: p}
			// originVal:    v,
			// 对每个节点进行跳点的直线可达性判断，并记录好跳点直线直线距离
			node.jumpDistance = searchStraightDis(g, jpGraph, p)
			nodes[y][x] = node
		}
	}

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			p := image.Pt(x, y)
			searchLeanDis(g, jpGraph, nodes, p)
		}
	}

	return &jpsPlusGraph{nodes: nodes}
}

func searchStraightDis(g JPSGraph, jpG jpGraph, p image.Point) [8]int {
	jumpDistance := [8]int{}
	if g.IsBlock(p) {
		return jumpDistance
	}
	for _, directionIdx := range straightDir {
		dir := directions[directionIdx]
		var distance int = 0
		for next := p.Add(dir); !g.IsBlock(next); next = next.Add(dir) {
			distance--
			if jpG[next.Y][next.X]&jumpDirections[directionIdx] != 0 {
				distance = -distance
				break
			}
		}
		jumpDistance[directionIdx] = distance
	}
	return jumpDistance
}

func searchLeanDis(g JPSGraph, jpG jpGraph, nodes [][]*node, p image.Point) {
	if g.IsBlock(p) {
		return
	}
	for _, directionIdx := range leanDir {
		dir := directions[directionIdx]
		// is near by block
		if g.IsBlock(p.Add(image.Pt(dir.X, 0))) || g.IsBlock(p.Add(image.Pt(0, dir.Y))) {
			continue
		}
		var distance int = 0
		for next := p.Add(dir); !g.IsBlock(next); next = next.Add(dir) {
			distance--
			if jpG[next.Y][next.X]&jumpDirections[directionIdx] != 0 {
				distance = -distance
				break
			}
			if nodes[next.Y][next.X].jumpDistance[leanToStraightDir[directionIdx][0]] > 0 {
				distance = -distance
				break
			}
			if nodes[next.Y][next.X].jumpDistance[leanToStraightDir[directionIdx][1]] > 0 {
				distance = -distance
				break
			}
		}
		nodes[p.Y][p.X].jumpDistance[directionIdx] = distance
	}
}
