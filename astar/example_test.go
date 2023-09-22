package astar

import (
	"fmt"
	"image"
	"math"
	"testing"
)

func TestExampleFindPath1(t *testing.T) {
	// Create a graph with 2D points as nodes
	p1 := image.Pt(3, 1)
	p2 := image.Pt(1, 2)
	p3 := image.Pt(2, 4)
	p4 := image.Pt(4, 5)
	p5 := image.Pt(4, 3)
	p6 := image.Pt(5, 1)
	p7 := image.Pt(8, 4)
	p8 := image.Pt(8, 3)
	p9 := image.Pt(6, 3)
	g := newGraph[image.Point]().
		link(p1, p2).link(p1, p3).
		link(p2, p3).link(p2, p4).
		link(p3, p4).link(p3, p5).
		link(p4, p6).link(p4, p7).
		link(p5, p7).
		link(p6, p9).
		link(p7, p8).
		link(p8, p9)

	// Find the shortest path from p1 to p9
	p := FindPath[image.Point](g, p1, p9, nodeDist, nodeDist)

	// Output the result
	if p == nil {
		fmt.Println("No path found.")
		return
	}
	for i, n := range p {
		fmt.Printf("%d: %s\n", i, n)
	}
	// Output:
	// 0: (3,1)
	// 1: (2,4)
	// 2: (4,5)
	// 3: (5,1)
	// 4: (6,3)
}

func nodeDist(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}

type graph[Node comparable] map[Node][]Node

func newGraph[Node comparable]() graph[Node] {
	return make(map[Node][]Node)
}

func (g graph[Node]) link(a, b Node) graph[Node] {
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
	return g
}

func (g graph[Node]) Neighbors(n, p, goal Node) []Node {
	return g[n]
}

func TestExampleFindPathMaze(t *testing.T) {
	maze := floorPlan{
		"###############",
		"#   # #     # #",
		"# ### ### ### #",
		"#   # # #   # #",
		"### # # # ### #",
		"# # #         #",
		"# # ### ### ###",
		"#   # # # #   #",
		"### # # # # ###",
		"# #       # # #",
		"# # ######### #",
		"#         #   #",
		"# ### # # ### #",
		"#   # # #     #",
		"###############",
	}
	start := image.Pt(1, 13) // Bottom left corner
	dest := image.Pt(13, 1)  // Top right corner

	// Find the shortest path
	path := FindPath[image.Point](maze, start, dest, distance, distance)

	for _, p := range path {
		maze.put(p, '.')
	}
	maze.print()
	// Output:
	// ###############
	// #   # #     #.#
	// # ### ### ###.#
	// #   # # #   #.#
	// ### # # # ###.#
	// # # #  .......#
	// # # ###.### ###
	// #   # #.# #   #
	// ### # #.# # ###
	// # #.....  # # #
	// # #.######### #
	// #...      #   #
	// #.### # # ### #
	// #.  # # #     #
	// ###############
}

func distance(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}

type floorPlan []string

func (f floorPlan) Neighbors(n, p, goal image.Point) []image.Point {
	offsets := []image.Point{
		image.Pt(0, -1), // North
		image.Pt(1, 0),  // East
		image.Pt(0, 1),  // South
		image.Pt(-1, 0), // West
	}
	res := make([]image.Point, 0, 4)
	for _, off := range offsets {
		q := n.Add(off)
		if f.isFreeAt(q) {
			res = append(res, q)
		}
	}
	return res
}

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
