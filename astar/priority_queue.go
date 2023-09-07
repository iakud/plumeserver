package astar

type item[T any] struct {
	value    T
	priority float64
}

type priorityQueue[T any] []*item[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
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
