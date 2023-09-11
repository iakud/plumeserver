package astar

import (
	"container/heap"
)

type Graph[Node any] interface {
	Neighbours(n Node) []Node
}

type CostFunc[Node any] func(a, b Node) float64

func FindPath[Node comparable](g Graph[Node], start, dest Node, d, h CostFunc[Node]) []Node {
	closeList := make(map[Node]struct{})
	openList := make(map[Node]*item[Node])

	var pq priorityQueue[Node]
	heap.Init(&pq)
	heap.Push(&pq, &item[Node]{value: start})

	for pq.Len() > 0 {
		it := heap.Pop(&pq).(*item[Node])
		node := it.value

		delete(openList, it.value)
		if node == dest {
			// Path found
			nodes := make([]Node, 0)
			for itLink := it; itLink != nil; itLink = itLink.link {
				nodes = append(nodes, itLink.value)
			}
			for i := 0; i < len(nodes); i++ {
				j := len(nodes) - i - 1
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
			return nodes
		}

		closeList[node] = struct{}{}

		for _, next := range g.Neighbours(node) {
			if _, ok := closeList[next]; ok {
				continue
			}
			gScore := d(node, next) + it.gScore
			itNext, ok := openList[next]
			if !ok {
				itNext = &item[Node]{value: next, link: it, gScore: gScore, fScore: gScore + h(next, dest)}
				openList[next] = itNext
				heap.Push(&pq, itNext)
				continue
			}
			if itNext.gScore <= gScore {
				continue
			}
			// update
			itNext.link = it
			itNext.gScore = gScore
			itNext.fScore = gScore + h(next, dest)
			heap.Fix(&pq, itNext.index)
		}
	}

	// No path found
	return nil
}
