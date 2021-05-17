package queue

import (
	"log"
	"sync"

	"github.com/rlaskowski/go-email/email"
)

var (
	EmailQueueType QueueType = "emailqueue"
)

type QueueType string

type Queue struct {
	queueProcess QueueProcess
}

type QueueFactory struct {
	factory map[QueueType]*Queue
	email   *email.Email
	mutex   *sync.Mutex
}

func NewFactory(email *email.Email) *QueueFactory {
	return &QueueFactory{
		factory: make(map[QueueType]*Queue),
		email:   email,
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
		if err := queue.queueProcess.Stop(); err != nil {
			log.Fatalf("Caught error while stopping queue %s", err.Error())
		}
	}

	return nil
}

func (q *QueueFactory) GetOrCreate(key QueueType) (QueueProcess, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if queue, ok := q.factory[key]; ok {
		return queue.queueProcess, nil
	}

	return q.createEmailQueue(key)
}

func (q *QueueFactory) createEmailQueue(key QueueType) (*EmailQueue, error) {
	emailq := NewEmailQueue(q.email)

	if err := emailq.Start(); err != nil {
		return nil, err
	}

	q.factory[key] = &Queue{
		queueProcess: emailq,
	}

	return emailq, nil
}
