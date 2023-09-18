package astar

import (
	"container/heap"
)

type Graph[Node any] interface {
	Neighbors(n Node) []Node
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

		for _, neighbor := range g.Neighbors(currentNode.node) {
			if _, ok := closeList[neighbor]; ok {
				continue
			}

			cost := d(currentNode.node, neighbor) + currentNode.gScore
			neighbourNode, ok := openList[neighbor]
			if !ok {
				// add
				neighbourNode = &item[Node]{
					node:   neighbor,
					gScore: cost,
					fScore: cost + h(neighbor, goal),
				}
				cameFrom[neighbor] = currentNode.node
				heap.Push(&pq, neighbourNode)
				openList[neighbor] = neighbourNode
				continue
			}
			// update
			if neighbourNode.gScore <= cost {
				continue
			}
			neighbourNode.gScore = cost
			neighbourNode.fScore = cost + h(neighbor, goal)
			heap.Fix(&pq, neighbourNode.index)
			cameFrom[neighbor] = currentNode.node
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
