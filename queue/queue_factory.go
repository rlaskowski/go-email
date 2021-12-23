package queue

import (
	"sync"
)

type Queue struct {
	queueProcess QueueProcess
}

type QueueFactory struct {
	factory map[string]*Queue
	mutex   *sync.Mutex
}

func NewFactory() *QueueFactory {
	return &QueueFactory{
		factory: make(map[string]*Queue),
		mutex:   &sync.Mutex{},
	}
}

func (q *QueueFactory) GetOrCreate(key string) QueueProcess {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if queue, ok := q.factory[key]; ok {
		return queue.queueProcess
	}

	return q.createPriorityQueue(key)
}

func (q *QueueFactory) createPriorityQueue(key string) *PriorityQueue {
	emailq := NewPriorityQueue()

	q.factory[key] = &Queue{
		queueProcess: emailq,
	}

	return emailq
}
