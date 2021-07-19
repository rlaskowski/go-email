package queue

import (
	"errors"
	"sync"

	"github.com/rlaskowski/go-email/email"
)

type QueueBox struct {
	emailPool sync.Pool
}

func (q *QueueBox) Start() error {
	q.emailPool.New = func() interface{} {
		return new(email.Email)
	}

	return q.start()
}

func (q *QueueBox) Stop() error {
	return nil
}

func (q *QueueBox) start() error {
	return errors.New("Not yet implemented")
}

func (q *QueueBox) acquireEmail() *email.Email {
	e := q.emailPool.Get().(*email.Email)
	defer q.emailPool.Put(e)

	return e
}
