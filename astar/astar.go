package astar

import (
	"container/heap"
)

type Graph[fScore any] interface {
	Neighbours(n fScore) []fScore
}

type CostFunc[Node any] func(a, b Node) float64

func FindPath[Node comparable](g Graph[Node], start, goal Node, d, h CostFunc[Node]) []Node {
	closeList := make(map[Node]struct{})
	openList := make(map[Node]*item[Node])

	// For node n, cameFrom[n] is the node immediately preceding it on the cheapest path from the start to n currently known.
	cameFrom := make(map[Node]Node)

	var pq priorityQueue[Node]
	heap.Init(&pq)
	// start
	startNode := &item[Node]{node: start}
	openList[start] = startNode
	heap.Push(&pq, startNode)

	for pq.Len() > 0 {
		currentNode := heap.Pop(&pq).(*item[Node])

		if currentNode.node == goal {
			// Path found
			return reconstructPath(cameFrom, goal)
		}

		delete(openList, currentNode.node)
		closeList[currentNode.node] = struct{}{}

		for _, neighbour := range g.Neighbours(currentNode.node) {
			if _, ok := closeList[neighbour]; ok {
				continue
			}

			cost := d(currentNode.node, neighbour) + currentNode.gScore
			neighbourNode, ok := openList[neighbour]
			if !ok {
				// add
				neighbourNode = &item[Node]{
					node:   neighbour,
					gScore: cost,
					fScore: cost + h(neighbour, goal),
				}
				cameFrom[neighbour] = currentNode.node
				heap.Push(&pq, neighbourNode)
				openList[neighbour] = neighbourNode
				continue
			}
			// update
			if neighbourNode.gScore <= cost {
				continue
			}
			neighbourNode.gScore = cost
			neighbourNode.fScore = cost + h(neighbour, goal)
			heap.Fix(&pq, neighbourNode.index)
			cameFrom[neighbour] = currentNode.node
		}
	}
	// No path found
	return nil
}

func reconstructPath[Node comparable](cameFrom map[Node]Node, current Node) []Node {
	var nodes = []Node{current}
	for node, ok := cameFrom[current]; ok; node, ok = cameFrom[node] {
		nodes = append(nodes, node)
	}
	length := len(nodes)
	reversedNodes := make([]Node, length)
	for i := 0; i < len(nodes); i++ {
		reversedNodes[i] = nodes[length-i-1]
	}
	return nodes
}
