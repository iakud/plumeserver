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

func FindPath[Node comparable](g Graph[Node], start, dest Node, d, h CostFunc[Node]) Path[Node] {
	closeList := make(map[Node]struct{})
	openList := make(map[Node]*item[Path[Node]])

	var pq priorityQueue[Path[Node]]
	heap.Init(&pq)
	heap.Push(&pq, &item[Path[Node]]{value: newPath(start)})

	for pq.Len() > 0 {
		p := heap.Pop(&pq).(*item[Path[Node]]).value
		n := p.last()
		if n == dest {
			// Path found
			return p
		}
		delete(openList, n)
		closeList[n] = struct{}{}

		for _, nb := range g.Neighbours(n) {
			if _, ok := closeList[nb]; ok {
				continue
			}
			cp := p.cont(nb)
			priority := cp.Cost(d) + h(nb, dest)
			itemNode, ok := openList[nb]
			if !ok {
				itemNode = &item[Path[Node]]{value: cp, priority: priority}
				openList[nb] = itemNode
				heap.Push(&pq, itemNode)
				continue
			}
			if itemNode.priority <= priority {
				continue
			}
			itemNode.value = cp
			itemNode.priority = priority
			heap.Fix(&pq, itemNode.index)
		}
	}

	// No path found
	return nil
}
