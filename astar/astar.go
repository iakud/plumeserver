package astar

import (
	"container/heap"
)

type Graph[Node any] interface {
	Neighbours(n Node) []Node
}

type CostFunc[Node any] func(a, b Node) float64

func FindPath[Node comparable](g Graph[Node], start, end Node, d, h CostFunc[Node]) []Node {
	closeList := make(map[Node]struct{})
	openList := make(map[Node]*item[Node])

	var pq priorityQueue[Node]
	heap.Init(&pq)
	// start
	from := &item[Node]{node: start}
	openList[start] = from
	heap.Push(&pq, from)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*item[Node])

		if current.node == end {
			// Path found
			return paths(current)
		}

		delete(openList, current.node)
		closeList[current.node] = struct{}{}

		for _, neighbour := range g.Neighbours(current.node) {
			if _, ok := closeList[neighbour]; ok {
				continue
			}

			cost := d(current.node, neighbour) + current.cost
			next, ok := openList[neighbour]
			if !ok {
				// add
				next = &item[Node]{
					node:     neighbour,
					from:     current,
					cost:     cost,
					priority: cost + h(neighbour, end),
				}
				openList[neighbour] = next
				heap.Push(&pq, next)
				continue
			}
			// update
			if next.cost <= cost {
				continue
			}
			next.from = current
			next.cost = cost
			next.priority = cost + h(neighbour, end)
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
