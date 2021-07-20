package queue

import (
	"errors"
	"sync"

	"github.com/rlaskowski/go-email/email"
)

type QueueBox struct {
	emailPool      sync.Pool
	receivingQueue QueueProcess
	sendingQueue   QueueProcess
}

func (q *QueueBox) Start() error {
	q.emailPool.New = func() interface{} {
		return email.NewEmail()
	}

	return q.start()
}

func (q *QueueBox) Stop() error {
	return nil
}

func (q *QueueBox) start() error {
	if err := q.initEmail(); err != nil {
		return err
	}

	go q.receiving()

	return nil
}

func (q *QueueBox) receiving() error {
	return errors.New("Not implemented yet")
}

func (q *QueueBox) sending() error {
	return errors.New("Not implemented yet")
}

func (q *QueueBox) initEmail() error {
	e := q.acquireEmail()
	return e.Init()
}

func (q *QueueBox) acquireEmail() *email.Email {
	e := q.emailPool.Get().(*email.Email)
	defer q.emailPool.Put(e)

	return e
}
