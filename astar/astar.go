package astar

import (
	"container/heap"
)

type Graph[Node any] interface {
	Neighbours(n Node) []Node
}

type CostFunc[Node any] func(a, b Node) float64
type Path[Node any] []Node

func newPath[Node any](start Node) Path[Node] {
	return []Node{start}
}

func (p Path[Node]) last() Node {
	return p[len(p)-1]
}

func (p Path[Node]) cont(n Node) Path[Node] {
	cp := make([]Node, len(p), len(p)+1)
	copy(cp, p)
	cp = append(cp, n)
	return cp
}

func (p Path[Node]) Cost(d CostFunc[Node]) (c float64) {
	for i := 1; i < len(p); i++ {
		c += d(p[i-1], p[i])
	}
	return c
}

func FindPath[Node comparable](g Graph[Node], start, dest Node, d, h CostFunc[Node]) []Node {
	closeList := make(map[Node]struct{})
	openList := make(map[Node]*item[Node])

	var pq priorityQueue[Node]
	heap.Init(&pq)
	heap.Push(&pq, &item[Node]{value: start})

	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*item[Node])
		n := node.value
		if n == dest {
			// Path found
			path := make([]Node, 0)
			for node != nil {
				path = append(path, node.value)
				node = node.from
			}
			for i := 0; i < len(path); i++ {
				j := len(path) - i - 1
				path[i], path[j] = path[j], path[i]
			}
			return path
		}
		delete(openList, n)
		closeList[n] = struct{}{}

		for _, nb := range g.Neighbours(n) {
			if _, ok := closeList[nb]; ok {
				continue
			}
			gScore := d(n, nb) + node.gScore
			itemNode, ok := openList[nb]
			if !ok {
				itemNode = &item[Node]{value: nb, from: node, gScore: gScore, fScore: gScore + h(nb, dest)}
				openList[nb] = itemNode
				heap.Push(&pq, itemNode)
				continue
			}
			if itemNode.gScore <= gScore {
				continue
			}
			itemNode.from = node
			itemNode.gScore = gScore
			itemNode.fScore = gScore + h(nb, dest)
			heap.Fix(&pq, itemNode.index)
		}
	}

	// No path found
	return nil
}
