package main

// import "container/heap"

type PriorityQueue []Job

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Job))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)

	job := old[n-1]
	*pq = old[:n-1]
	return job
}
