package queue

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/model"
)

type EmailQueue struct {
	mutex *sync.Mutex
	queue []interface{}
}

func NewEmailQueue() *EmailQueue {
	return &EmailQueue{
		mutex: &sync.Mutex{},
	}
}

func (e *EmailQueue) Start() error {
	go e.scan()
	return nil
}

func (e *EmailQueue) Stop() error {
	return nil
}

func (e *EmailQueue) Publish(message interface{}) error {
	e.queue = append(e.queue, message)
	return nil
}

func (e *EmailQueue) Subscribe() error {
	firstElement, err := e.firstElement()
	if err != nil {
		return err
	}

	message := firstElement.(*model.Message)

	email := email.NewEmail()

	go email.Send(message.Key, message)

	return nil
}

func (e *EmailQueue) isEmpty() bool {
	if !(len(e.queue) > 0) {
		return true
	}
	return false
}

func (e *EmailQueue) lastElement() (interface{}, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.isEmpty() {
		return nil, errors.New("No elements")
	}

	element := e.queue[len(e.queue)]
	e.queue = e.queue[:len(e.queue)-1]

	return element, nil
}

func (e *EmailQueue) firstElement() (interface{}, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.isEmpty() {
		return nil, errors.New("No elements")
	}

	element := e.queue[0]
	e.queue = e.queue[1:]

	return element, nil
}

func (e *EmailQueue) scan() {
	for {

		if err := e.Subscribe(); err != nil {
			log.Printf("Could not get element from queue: %s", err)
		}

		time.Sleep(time.Second * 3)
	}
}
