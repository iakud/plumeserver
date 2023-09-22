package jpsplus

import (
	"fmt"
	"image"
	"log"
	"math"
	"testing"

	"github.com/iakud/plumeserver/astar"
)

type grid []string

func (g grid) IsFreeAt(p image.Point) bool {
	return !g.IsOutsideAt(p) && !g.IsObstacleAt(p)
}

func (g grid) IsOutsideAt(p image.Point) bool {
	return p.Y < 0 || p.Y >= len(g) || p.X < 0 || p.X >= len(g[p.Y])
}

func (g grid) IsObstacleAt(p image.Point) bool {
	return g[p.Y][p.X] != ' '
}

func (g grid) put(p image.Point, c rune) {
	g[p.Y] = g[p.Y][:p.X] + string(c) + g[p.Y][p.X+1:]
}

func (g grid) print() {
	for _, y := range g {
		fmt.Println(y)
	}
}

func TestJump(t *testing.T) {
	var g = grid{
		" ##                         ##                    ",
		"                            ##                    ",
		"       ###         ##       ##                    ",
		"       ###         ##       ##                    ",
		"       ###         ##       ##                    ",
		"       ###         ##       ##                    ",
		"                   ##       ##                    ",
		"                   ###########                    ",
		"             ########                             ",
		"             ########                             ",
		"                                                  ",
		"                                                  ",
		"                            ##                    ",
		"                 #          ##                    ",
		"                 #          ##     ###############",
		"                 #          ##                    ",
		"                 #          ##                    ",
		"##################          ##                    ",
		"                            ##                    ",
		"                            ##############        ",
		"                                                  ",
		"                                                  ",
		"                 #############                    ",
		"                            ##                    ",
		"       ###                  ##                    ",
		"       ###                  ##                    ",
		"       ###                  ##                    ",
		"       ###                                        ",
		"       ###                                        ",
		"       ###                                        ",
		"       ###                ###################     ",
		"       ###                ##            ##        ",
		"       ###                ##            ##        ",
		"       ###                ##            ##        ",
		"       ###                ##            ##        ",
		"################          ##            ##        ",
		"                          ##            ##        ",
		"                          ##            ##        ",
		"                          ##            ##        ",
		"                          ##            ##        ",
		"                          ##            ##        ",
		"              ##############            ##        ",
		"                                        ##        ",
		"                           ###############        ",
		"                           ##                     ",
		"          ###              ##             ####    ",
		"          ###              ##             ##      ",
		"          ###              ##             ##      ",
		"          ###                             ##      ",
		"          ###                             ##      ",
	}

	start := image.Pt(0, 0)
	goal := image.Pt(49, 49)
	graph := NewGraph(g, image.Pt(len(g), len(g[0])))

	path := astar.FindPath(graph, start, goal, distance, distance)
	log.Println("path:", len(path))

	for _, p := range path {
		g.put(p, '.')
	}
	g.print()
}

func distance(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}
