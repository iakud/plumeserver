package astar

import (
	"container/heap"
)

type Graph[Node any] interface {
	Neighbours(n Node) []Node
}

type CostFunc[Node any] func(a, b Node) float64

func FindPath[Node comparable](g Graph[Node], start, end Node, d, h CostFunc[Node]) []Node {
	nodeList := make(map[Node]*item[Node])

	var pq priorityQueue[Node]
	heap.Init(&pq)
	// start item
	startItem := &item[Node]{node: start}
	nodeList[start] = startItem
	heap.Push(&pq, startItem)

	for pq.Len() > 0 {
		currentItem := heap.Pop(&pq).(*item[Node])

		if currentItem.node == end {
			// Path found
			return paths(currentItem)
		}

		currentItem.closed = true // close

		for _, neighbour := range g.Neighbours(currentItem.node) {
			neighbourItem, ok := nodeList[neighbour]
			if !ok {
				// new item
				cost := d(currentItem.node, neighbour) + currentItem.cost
				neighbourItem = &item[Node]{node: neighbour, from: currentItem, cost: cost, priority: cost + h(neighbour, end)}
				nodeList[neighbour] = neighbourItem
				heap.Push(&pq, neighbourItem)
				continue
			}
			if neighbourItem.closed {
				continue
			}
			// update item
			cost := d(currentItem.node, neighbour) + currentItem.cost
			if neighbourItem.cost <= cost {
				continue
			}
			neighbourItem.from = currentItem
			neighbourItem.cost = cost
			neighbourItem.priority = cost + h(neighbour, end)
			heap.Fix(&pq, neighbourItem.index)
		}
	}
	// No path found
	return nil
}

func paths[Node any](endItem *item[Node]) []Node {
	nodes := make([]Node, 0)
	for from := endItem; from != nil; from = from.from {
		nodes = append(nodes, from.node)
	}
	for i := 0; i < len(nodes); i++ {
		j := len(nodes) - i - 1
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes
}
