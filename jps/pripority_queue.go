package jps

type item[T any] struct {
	node      T
	from      *item[T]
	direction DirectionIdx
	cost      float64 // G
	priority  float64 // F (F = G + H)
	index     int
}

type priorityQueue[T any] []*item[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *priorityQueue[T]) Push(x any) {
	it := x.(*item[T])
	it.index = len(*pq)
	*pq = append(*pq, it)
}

func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[0 : n-1]
	return it
}
