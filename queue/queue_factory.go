package queue

import (
	"log"
	"sync"
)

var (
	EmailQueueType QueueType = "emailqueue"
)

type QueueType string

type Queue struct {
	queueConnection QueueConnection
}

type QueueFactory struct {
	factory map[QueueType]Queue
	mutex   *sync.Mutex
}

func NewFactory() *QueueFactory {
	return &QueueFactory{
		factory: make(map[QueueType]Queue),
		mutex:   &sync.Mutex{},
	}
}

func (q *QueueFactory) Start() error {
	log.Print("Starting Queue Factory")
	return nil
}

func (q *QueueFactory) Stop() error {
	log.Print("Stopping Queue Factory")

	for _, queue := range q.factory {
		if err := queue.queueConnection.Stop(); err != nil {
			log.Fatalf("Caught error while stopping queue %s", err.Error())
		}
	}

	return nil
}

func (q *QueueFactory) GetOrCreate(key QueueType) (QueueConnection, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if queue, ok := q.factory[key]; ok {
		return queue.queueConnection, nil
	}

	return q.createEmailQueue(key)
}

func (q *QueueFactory) createEmailQueue(key QueueType) (*EmailQueue, error) {
	emailq := NewEmailQueue()

	if err := emailq.Start(); err != nil {
		return nil, err
	}

	q.factory[key] = Queue{
		queueConnection: emailq,
	}

	return emailq, nil
}
