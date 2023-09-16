package jps

import (
	"container/heap"
	"image"
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
)

type Graph[Node any] interface {
	Neighbours(n Node) []Node
}

type CostFunc[Node any] func(a, b Node) float64

func FindPath(g *jpsPlusGraph, start, end image.Point, d, h CostFunc[image.Point]) []image.Point {
	closeList := make(map[image.Point]struct{})
	openList := make(map[image.Point]*item[image.Point])

	var pq priorityQueue[image.Point]
	heap.Init(&pq)
	// start
	from := &item[image.Point]{node: start, direction: IdxAll}
	openList[start] = from
	heap.Push(&pq, from)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*item[image.Point])

		if current.node == end {
			// Path found
			return paths(current)
		}

		delete(openList, current.node)
		closeList[current.node] = struct{}{}

		jumpDistance := g.nodes[current.node.Y][current.node.X].jumpDistance

		for _, direction := range directionMap[current.direction] {
			distance := jumpDistance[direction]
			if distance == 0 {
				continue
			}
			dir := directions[direction]
			if distance < 0 {
				distance = -distance
			}
			to := current.node.Add(dir.Mul(distance))
			endDir := end.Sub(current.node)
			// 方向平行
			if endDir.Y*dir.X == endDir.X*dir.Y && current.node.X <= end.X && end.X <= to.X && current.node.Y <= end.Y && end.Y <= to.Y {
				to = end
			}

			// open
			cost := d(current.node, to) + current.cost
			next, ok := openList[to]
			if !ok {
				// add
				next = &item[image.Point]{
					node:      to,
					from:      current,
					direction: direction,
					cost:      cost,
					priority:  cost + h(to, end),
				}
				openList[to] = next
				heap.Push(&pq, next)
				continue
			}
			// update
			if next.cost <= cost {
				continue
			}
			next.from = current
			next.direction = direction
			next.cost = cost
			next.priority = cost + h(to, end)
			heap.Fix(&pq, next.index)
		}
	}
	// No path found
	return nil
}

func paths[Node any](current *item[Node]) []Node {
	nodes := make([]Node, 0)
	for from := current; from != nil; from = from.from {
		nodes = append(nodes, from.node)
	}
	for i := 0; i < len(nodes); i++ {
		j := len(nodes) - i - 1
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes
}
