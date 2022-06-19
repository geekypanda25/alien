package queue

import "container/heap"

type (
	Heapable interface {
		Priority(other interface{}) bool
	}

	items []Heapable

	PriorityQueue struct {
		queue *items
	}
)

// NewPriorityQueue references a new queue
func NewPriorityQueue() *PriorityQueue {
	
	pq := &PriorityQueue{queue: &items{}}

	heap.Init(pq.queue)
	
	return pq
}

// Adds element to queue
func (pq *PriorityQueue) Push(s Heapable) {
	
	heap.Push(pq.queue, s)

}

// Removes item from queue
func (pq *PriorityQueue) Pop() Heapable {

	return heap.Pop(pq.queue).(Heapable)

}

// Size of the queue.
func (pq *PriorityQueue) Size() int {

	return pq.queue.Len()

}

// Implements the sort interface
func (pq items) Len() int {

	return len(pq)

}

// Implements the sort interface
func (pq items) Less(i, j int) bool {

	return pq[i].Priority(pq[j])

}

// Implements the sort interface
func (pq items) Swap(i, j int) {

	pq[i], pq[j] = pq[j], pq[i]

}

// Implements the sort interface
func (pq *items) Push(x interface{}) {

	item := x.(Heapable)

	*pq = append(*pq, item)

}

// Implements the sort interface
func (pq *items) Pop() interface{} {

	old := *pq

	n := len(old)

	item := old[n-1]

	*pq = old[0 : n-1]

	return item
}
