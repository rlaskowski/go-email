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
	mutex     *sync.Mutex
	queue     []QueueStore
	emailPool sync.Pool
	buff      bytes.Buffer
}

type QueueStore struct {
	Subject QueueSubject
	Message interface{}
}

func NewEmailQueue() *EmailQueue {
	e := &EmailQueue{
		mutex: &sync.Mutex{},
		queue: make([]QueueStore, 0),
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
	e.queue = []QueueStore{}
	return nil
}

func (e *EmailQueue) Publish(Subject QueueSubject, message ...interface{}) error {
	var storeMessage QueueStore

	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.prepareData(message); err != nil {
		return err
	}

	storeMessage = QueueStore{Subject: Subject, Message: message}

	e.queue = append(e.queue, storeMessage)

	return nil
}

func (e *EmailQueue) prepareData(data []interface{}) error {
	for i, m := range data {
		t := reflect.TypeOf(m).Elem()

		if t == reflect.TypeOf(&model.File{}).Elem() {
			f := m.(*model.File)

			if err := e.storeOrRead(f); err != nil {
				return err
			}

			data[i] = f
		}
	}
	return nil
}

func (e *EmailQueue) Subscribe(Subject QueueSubject) error {
	return nil
}

func (e *EmailQueue) send() error {
	var message *model.Message
	var file *model.File

	firstElement, err := e.firstElement()
	if err != nil {
		return err
	}

	for _, m := range firstElement.([]interface{}) {
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
		em := e.acquireEmail()

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
	em := e.emailPool.Get().(*email.Email)
	defer e.emailPool.Put(em)

	return em
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

func (e *EmailQueue) receive() error {
	em := e.acquireEmail()
	em.Receive(func(info email.Stat) error {
		e.Publish(SubjectReceiving, info)

		return nil
	})

	return nil
}

func (e *EmailQueue) scan() {
	for {

		if err := e.send(); err != nil {
			time.Sleep(time.Second * 3)
		}

	}
}
