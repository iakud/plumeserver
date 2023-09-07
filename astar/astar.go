package astar

import "container/heap"

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
	closed := make(map[Node]bool)

	var pq priorityQueue[Path[Node]]
	heap.Init(&pq)
	heap.Push(&pq, &item[Path[Node]]{value: newPath(start)})

	for pq.Len() > 0 {
		p := heap.Pop(&pq).(*item[Path[Node]]).value
		n := p.last()
		if closed[n] {
			continue
		}
		if n == dest {
			// Path found
			return p
		}
		closed[n] = true

		for _, nb := range g.Neighbours(n) {
			cp := p.cont(nb)
			heap.Push(&pq, &item[Path[Node]]{
				value:    cp,
				priority: cp.Cost(d) + h(nb, dest),
			})
		}
	}

	// No path found
	return nil
}
