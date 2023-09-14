package jps

import (
	"fmt"
	"image"

	"golang.org/x/image/vector"
)

type node struct {
	image.Point
	originVal    uint8
	jumpDistance [8]int
}

const (
	JPSPlusUp = 1 << iota
	JPSPlusRight
	JPSPlusDown
	JPSPlusLeft
)

var (
	UpLeft    = image.Point{X: -1, Y: -1}
	Up        = image.Point{X: 0, Y: -1}
	UpRight   = image.Point{X: 1, Y: -1}
	Right     = image.Point{X: 1, Y: 0}
	DownRight = image.Point{X: 1, Y: 1}
	Down      = image.Point{X: 0, Y: 1}
	DownLeft  = image.Point{X: -1, Y: 1}
	Left      = image.Point{X: -1, Y: 0}
)

var (
	jpsDirections = [4]uint8{JPSPlusDown, JPSPlusUp, JPSPlusRight, JPSPlusLeft}
	jpsVectors    = [4]image.Point{Down, Up, Left, Right}
)

type JumpPointDefine struct {
	direction uint8
	vector    image.Point
}

func isJumpPoint(f floorPlan, p image.Point, dir image.Point) bool {
	if !f.isFreeAt(p.Sub(dir)) {
		return false
	}
	// forced neighbour
	if f.isFreeAt(image.Pt(p.X+dir.Y, p.Y+dir.X)) && !f.isFreeAt(image.Pt(p.X+dir.Y, p.Y+dir.X).Sub(dir)) {
		return true
	}
	if f.isFreeAt(image.Pt(p.X-dir.Y, p.Y-dir.X)) && !f.isFreeAt(image.Pt(p.X-dir.Y, p.Y-dir.X).Sub(dir)) {
		return true
	}
	return false
}

var JumpPointDir = [4]JumpPointDefine{
	JumpPointDefine{direction: JPSPlusDown, vector: Down},
	JumpPointDefine{direction: JPSPlusUp, vector: Up},
	JumpPointDefine{direction: JPSPlusRight, vector: Left},
	JumpPointDefine{direction: JPSPlusLeft, vector: Right},
}

type floorPlan []string

func (f floorPlan) isFreeAt(p image.Point) bool {
	return f.isInBounds(p) && f[p.Y][p.X] == ' '
}

func (f floorPlan) isInBounds(p image.Point) bool {
	return (0 <= p.X && p.X < len(f[p.Y])) && (0 <= p.Y && p.Y < len(f))
}

func (f floorPlan) put(p image.Point, c rune) {
	f[p.Y] = f[p.Y][:p.X] + string(c) + f[p.Y][p.X+1:]
}

func (f floorPlan) print() {
	for _, row := range f {
		fmt.Println(row)
	}
}

type jpGraph [][]uint8

func newJumpPointGraph(f floorPlan) jpGraph {
	jpGraph := make([][]uint8, len(f))
	for y := range f {
		jpGraph[y] = make([]uint8, len(f[y]))
		for x := range f[y] {
			if f[y][x] != ' ' {
				continue
			}
			p := image.Pt(x, y)
			var direction uint8
			for index, dir := range jpsVectors {
				if !isJumpPoint(f, p, dir) {
					continue
				}
				direction |= jpsDirections[index]
			}
			jpGraph[y][x] = direction
		}
	}
	return jpGraph
}

type jpsplusGraph struct {
	floorPlan
	nodes [][]*node
}

func preCptJpMatrix(f floorPlan) *jpsplusGraph {
	jpGraph := newJumpPointGraph(f)
	nodes := make([][]*node, len(f))
	for y := range f {
		nodes[y] = make([]*node, len(f[y]))
		for x, _ := range f[y] {
			p := image.Pt(x, y)
			node := &node{Point: p}
			// originVal:    v,
			// 对每个节点进行跳点的直线可达性判断，并记录好跳点直线直线距离
			node.jumpDistance = searchStraightDis(f, jpGraph, p)
			nodes[y][x] = node
		}
	}
	return &jpsplusGraph{floorPlan: f, nodes: nil}
}

func searchStraightDis(f floorPlan, jpG jpGraph, p image.Point) [8]int {

}
