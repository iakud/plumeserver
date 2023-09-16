package jps

import (
	"fmt"
	"image"
	"log"
	"math"
	"testing"
)

type graph []string

func (g graph) IsBlock(p image.Point) bool {
	return p.Y < 0 || p.Y >= len(g) || p.X < 0 || p.X >= len(g[p.Y]) || g[p.Y][p.X] != ' '
}

func (g graph) Size() image.Point {
	return image.Pt(len(g), len(g[0]))
}

func (g graph) put(p image.Point, c rune) {
	g[p.Y] = g[p.Y][:p.X] + string(c) + g[p.Y][p.X+1:]
}

func (g graph) print() {
	for _, y := range g {
		fmt.Println(y)
	}
}

func TestJump(t *testing.T) {
	var g = graph{
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
	end := image.Pt(49, 49)
	// jp := newJumpPointGraph(g)
	graph := newJpsPlusGraph(g)

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			graph.nodes[start.Y][start.X].printDis()
		}
	}

	path := FindPath(graph, start, end, distance, distance)
	log.Println("path:", len(path))
	g.put(start, 'S')
	g.put(end, 'E')
	for _, p := range path {
		g.put(p, 'o')
	}
	g.print()
}

func distance(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}
