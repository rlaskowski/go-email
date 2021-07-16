package queue

import (
	"sync"

	"github.com/rlaskowski/go-email/email"
)

type Queue struct {
	queueProcess QueueProcess
}

type QueueFactory struct {
	factory map[string]*Queue
	email   *email.Email
	mutex   *sync.Mutex
}

func NewFactory(email *email.Email) *QueueFactory {
	return &QueueFactory{
		factory: make(map[string]*Queue),
		email:   email,
		mutex:   &sync.Mutex{},
	}
}

func (q *QueueFactory) GetOrCreate(key string) (QueueProcess, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if queue, ok := q.factory[key]; ok {
		return queue.queueProcess, nil
	}

	return q.createPriorityQueue(key)
}

func (q *QueueFactory) createPriorityQueue(key string) (*PriorityQueue, error) {
	emailq := NewPriorityQueue()

	q.factory[key] = &Queue{
		queueProcess: emailq,
	}

	return emailq, nil
}
