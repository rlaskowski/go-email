package queue

import (
	"errors"

	"github.com/rlaskowski/go-email"
)

type Queue struct {
	message []*email.Message
}

type QueueSender interface {
	Publish(message *email.Message) error
	Subscribe() ([]*email.Message, error)
}

func NewQueue() *Queue {
	return &Queue{
		message: make([]*email.Message, 0),
	}
}

func (q *Queue) isEmpty() bool {
	if !(len(q.message) > 0) {
		return true
	}
	return false
}

func (q *Queue) Publish(message *email.Message) error {
	if q.message == nil {
		return errors.New("Queue is not initialized")
	}
	q.message = append(q.message, message)

	return nil
}

func (q *Queue) Subscribe() ([]*email.Message, error) {
	if q.isEmpty() {
		return nil, errors.New("Queue is empty")
	}

	return q.message[1:], nil
}
