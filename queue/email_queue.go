package queue

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/model"
	"github.com/rlaskowski/go-email/store"
)

const (
	SubjectReceiving QueueSubject = "receiving"
	SubjectSending   QueueSubject = "sending"
)

type QueueSubject string

type EmailQueue struct {
	mutex *sync.Mutex
	queue []QueueStore
	email *email.Email
	buff  bytes.Buffer
}

type QueueStore struct {
	Subject QueueSubject
	Message []interface{}
}

func NewEmailQueue(email *email.Email) *EmailQueue {
	return &EmailQueue{
		mutex: &sync.Mutex{},
		queue: make([]QueueStore, 0),
		email: email,
	}
}

func (e *EmailQueue) Start() error {
	go e.start()
	return nil
}

func (e *EmailQueue) Stop() error {
	return nil
}

func (e *EmailQueue) Publish(Subject QueueSubject, message ...interface{}) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.prepareData(message); err != nil {
		return err
	}

	storeMessage := QueueStore{Subject: Subject, Message: message}
	e.enqueue(storeMessage)

	return nil
}

func (e *EmailQueue) Subscribe(Subject QueueSubject) ([]QueueStore, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var list []QueueStore
	if !(e.len() > 0) {
		return nil, errors.New("Not email to receive")
	}

	qslist := e.queue
	for _, l := range qslist {
		if l.Subject == SubjectReceiving {
			list = append(list, l)
		}
	}

	return list, nil
}

func (e *EmailQueue) prepareData(data []interface{}) error {
	for i, m := range data {
		if reflect.TypeOf(m).Kind() == reflect.Slice {
			t := reflect.TypeOf(m).Elem()

			if t == reflect.TypeOf(&model.File{}).Elem() {
				f := m.(*model.File)

				if err := e.storeOrRead(f); err != nil {
					return err
				}

				data[i] = f
			}
		}

	}
	return nil
}

func (e *EmailQueue) send() error {
	var message *model.Message
	var file *model.File

	dequeue := <-e.dequeue()

	if !(dequeue.Subject == SubjectSending) {
		return nil
	}

	for _, m := range dequeue.Message {
		t := reflect.TypeOf(m).Elem()

		switch t {
		case reflect.TypeOf(&model.Message{}).Elem():
			message = m.(*model.Message)
		case reflect.TypeOf(&model.File{}).Elem():
			file = m.(*model.File)
		}
	}

	select {
	case ret := <-e.sendEmail(message, file):
		if !ret {

			if err := e.Publish(SubjectSending, message, file); err != nil {
				return errors.New(err.Error())
			}
		}
	}
	return nil
}

func (e *EmailQueue) sendEmail(message *model.Message, file *model.File) <-chan bool {
	rCh := make(chan bool, 100)

	go func() {
		em := e.email

		err := em.Send(message, file)

		if err != nil {
			rCh <- false
			close(rCh)
			return
		}

		rCh <- true
		close(rCh)

	}()

	return rCh
}

func (e *EmailQueue) storeOrRead(file *model.File) error {
	if config.FileStore {

		uuid, err := e.storeFile(file)
		if err != nil {
			return err
		}

		file.Key = uuid

	} else {
		data, err := ioutil.ReadAll(file.Reader)
		if err != nil {
			return err
		}

		file.Data = data

	}

	return nil
}

func (e *EmailQueue) readFile(file []byte) ([]byte, error) {
	store := store.NewFileStore(config.FileStorePath)

	dir := store.ControllDir(string(file))
	filepath := filepath.Join(dir, string(file))

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (e *EmailQueue) storeFile(file *model.File) (string, error) {
	store := store.NewFileStore(config.FileStorePath)
	return store.Store(file.Reader)
}

func (e *EmailQueue) isEmpty() bool {
	if len(e.queue) == 0 {
		return true
	}
	return false
}

func (e *EmailQueue) len() int {
	return len(e.queue)
}

func (e *EmailQueue) enqueue(qstore QueueStore) {
	e.queue = append(e.queue, qstore)
}

func (e *EmailQueue) dequeue() <-chan QueueStore {
	qsch := make(chan QueueStore)

	go func() {

		if !e.isEmpty() {

			element := e.queue[0]
			e.queue = e.queue[1:]
			qsch <- element
		}

		close(qsch)

	}()

	return qsch
}

func (e *EmailQueue) receive() error {
	go func() {
		em := e.email
		em.Receive(func(info email.Stat) error {

			if !e.findReceive(info) {
				e.Publish(SubjectReceiving, info)
			}

			return nil
		})
	}()

	return nil
}

func (e *EmailQueue) findReceive(info email.Stat) bool {
	for _, i := range e.queue {
		if i.Subject == SubjectReceiving {
			if i.Message[0] == info {
				return true
			}
		}
	}
	return false
}

func (e *EmailQueue) start() {

	if err := e.send(); err != nil {

	}

	if err := e.receive(); err != nil {

	}

	time.Sleep(config.QueueRefreshTime)

}
