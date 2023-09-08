package astar

type item[T any] struct {
	value  T
	from   *item[T]
	gScore float64
	fScore float64
	index  int
}

type priorityQueue[T any] []*item[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].fScore < pq[j].fScore
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *priorityQueue[T]) Push(x any) {
	*pq = append(*pq, x.(*item[T]))
}

func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[0 : n-1]
	return it
}

type queue[Node comparable] struct {
	nodes         map[Node]*item[Node]
	priorityQueue []*item[Node]
}

func newQueue[Node comparable]() *queue[Node] {
	return &queue[Node]{nodes: make(map[Node]*item[Node])}
}

func (q *queue[Node]) Push(node Node) {
	// nodes
	// heap.Push(&q.priorityQueue, &item[Node]{value: node})
}
