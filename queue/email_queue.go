package queue

import (
	"bytes"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/model"
)

type EmailQueue struct {
	mutex     *sync.Mutex
	queue     []interface{}
	emailPool sync.Pool
	buff      bytes.Buffer
}

func NewEmailQueue() *EmailQueue {
	e := &EmailQueue{
		mutex: &sync.Mutex{},
	}

	e.emailPool.New = func() interface{} {
		return email.NewEmail()
	}

	return e
}

func (e *EmailQueue) Start() error {
	go e.scan()
	return nil
}

func (e *EmailQueue) Stop() error {
	return nil
}

func (e *EmailQueue) Publish(message ...interface{}) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.queue = append(e.queue, message)

	return nil
}

func (e *EmailQueue) Subscribe() error {
	return nil
}

func (e *EmailQueue) send() error {
	var message *model.Message
	var file string

	firstElement, err := e.firstElement()
	if err != nil {
		return err
	}

	v := reflect.ValueOf(firstElement)

	if v.Kind() == reflect.Slice {
		for _, m := range firstElement.([]interface{}) {
			switch v := reflect.ValueOf(m); v.Kind() {
			case reflect.String:
				file = m.(string)
			default:
				message = m.(*model.Message)
			}
		}
	}

	select {
	case ret := <-e.sendEmail(message, file):
		if !ret {

			if err := e.Publish(firstElement); err != nil {
				return errors.New(err.Error())
			}
			return errors.New("Error when try to send email")

		}
	}

	return nil
}

func (e *EmailQueue) sendEmail(message *model.Message, file string) <-chan bool {
	rCh := make(chan bool, 100)

	go func() {
		email := e.acquireEmail()
		err := email.Send(message, file)

		if err == nil {
			rCh <- true
		}

		close(rCh)
	}()

	return rCh
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

func (e *EmailQueue) acquireEmail() *email.Email {
	email := e.emailPool.Get().(*email.Email)
	defer e.emailPool.Put(email)

	return email
}

/* func (e *EmailQueue) WriteToTable() string {
	writer := new(tabwriter.Writer)

	e.buff.Reset()

	sink := bufio.NewWriter(&e.buff)
	writer.Init(sink, 0, 8, 2, ' ', tabwriter.Debug)

	_, _ = fmt.Fprintln(writer, "Host \t Pathname ")

	_ = d.IterateAll(func(database Database) error {
		_, _ = fmt.Fprintf(writer, "%s \t %s \n", database.databaseConnection.Host(), database.databaseConnection.Path())
		return nil
	})

	_, _ = fmt.Fprintln(writer)

	if err := writer.Flush(); err != nil {
		return err.Error()
	}

	if err := sink.Flush(); err != nil {
		return err.Error()
	}

	return e.buff.String()
} */

func (e *EmailQueue) scan() {
	for {

		if err := e.send(); err != nil {
			time.Sleep(time.Second * 3)
		}

	}
}
