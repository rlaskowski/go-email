package queue

import (
	"container/heap"

	"github.com/rlaskowski/go-email/email"
)

type QueueSubject string

type QueueStore struct {
	Message  interface{}
	Priority int
	Index    int
}

type PriorityQueue struct {
	queue []*QueueStore
	email *email.Email
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		queue: make([]*QueueStore, 0),
	}

	heap.Init(pq)

	return pq
}

func (p *PriorityQueue) Push(qstore interface{}) {
	n := len(p.queue)
	qs := qstore.(*QueueStore)
	qs.Index = n
	p.queue = append(p.queue, qs)
}

func (p *PriorityQueue) Pop() interface{} {
	old := *&p.queue
	n := len(old)
	qstore := old[n-1]
	//old[n-1] = nil
	qstore.Index = -1
	*&p.queue = old[0 : n-1]

	return qstore
}

func (p *PriorityQueue) Len() int {
	return len(p.queue)
}

func (p *PriorityQueue) Less(i, j int) bool {
	return p.queue[i].Priority > p.queue[j].Priority
}

func (p *PriorityQueue) Swap(i, j int) {
	p.queue[i], p.queue[j] = p.queue[j], p.queue[i]
	p.queue[i].Index = i
	p.queue[j].Index = j
}
