package queue

import (
	"container/heap"

	"github.com/rlaskowski/go-email/email"
)

const (
	SubjectReceiving QueueSubject = "receiving"
	SubjectSending   QueueSubject = "sending"
)

type QueueSubject string

type QueueStore struct {
	Subject  QueueSubject
	Message  interface{}
	Priority int
	Index    int
}

type EmailQueue struct {
	queue []*QueueStore
	email *email.Email
}

func NewEmailQueue(email *email.Email) *EmailQueue {
	return &EmailQueue{
		email: email,
		queue: make([]*QueueStore, 0),
	}
}

func (e *EmailQueue) Start() error {
	heap.Init(e)

	return nil
}

func (e *EmailQueue) Stop() error {
	return nil
}

func (e *EmailQueue) Push(qstore interface{}) {
	n := len(e.queue)
	qs := qstore.(*QueueStore)
	qs.Index = n
	e.queue = append(e.queue, qs)
}

func (e *EmailQueue) Pop() interface{} {
	old := *&e.queue
	n := len(old)
	qstore := old[n-1]
	//old[n-1] = nil
	qstore.Index = -1
	*&e.queue = old[0 : n-1]

	return qstore
}

func (e *EmailQueue) Len() int {
	return len(e.queue)
}

func (e *EmailQueue) Less(i, j int) bool {
	return e.queue[i].Priority > e.queue[j].Priority
}

func (e *EmailQueue) Swap(i, j int) {
	e.queue[i], e.queue[j] = e.queue[j], e.queue[i]
	e.queue[i].Index = i
	e.queue[j].Index = j
}

func (e *EmailQueue) findReceive(info email.Stat) bool {
	for _, i := range e.queue {
		if i.Subject == SubjectReceiving {
			if i.Message == info {
				return true
			}
		}
	}
	return false
}
