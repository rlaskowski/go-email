package queue

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/email"
)

const (
	Q_RECV string = "receiving"
	Q_SEND string = "sending"
)

type QueueBox struct {
	emailPool      sync.Pool
	queueFactory   *QueueFactory
	receivingQueue QueueProcess
	sendingQueue   QueueProcess
	serviceConfig  config.ServiceConfig
}

func NewQueuBox(serviceConfig config.ServiceConfig) *QueueBox {
	q := &QueueBox{
		queueFactory:  NewFactory(),
		serviceConfig: serviceConfig,
	}

	q.emailPool.New = func() interface{} {
		return email.NewEmail()
	}

	return q
}

func (q *QueueBox) Start() error {
	go q.receiving()

	go q.sending()

	return nil
}

func (q *QueueBox) Stop() error {
	return nil
}

func (q *QueueBox) receiving() {
	for {
		if err := q.receiveEmail(); err != nil {
			log.Printf("Couldn't read message due to: %s", err)
		}

		time.Sleep(q.serviceConfig.QueueRefreshTime)
	}

}

func (q *QueueBox) sending() {
	for {

		time.Sleep(q.serviceConfig.QueueRefreshTime)
	}
}

func (q *QueueBox) receiveEmail() error {
	e, err := q.acquireEmail()
	if err != nil {
		return err
	}

	for _, c := range e.Config() {
		sl, err := e.Stat(c.Key)
		if err != nil {
			return err
		}

		for _, s := range sl {
			mi, err := e.ReadMessage(c.Key, s.MessageNumber)
			if err != nil {
				return err
			}

			qid, err := q.queueId(c.Key, Q_RECV)
			if err != nil {
				return err
			}

			q.pushToQueue(qid, mi)
		}
	}

	return nil
}

func (q *QueueBox) ReceiveMessage(key string) ([]*email.MessageInfo, error) {
	qid, err := q.queueId(key, Q_RECV)
	if err != nil {
		return nil, err
	}

	pq := q.queueFactory.GetOrCreate(qid)

	//len := pq.Len()

	list := make([]*email.MessageInfo, 0)

	for pq.Len() > 0 {
		ps, ok := heap.Pop(pq).(*QueueStore)
		if !ok {
			return nil, errors.New("Error in pull queue, bad parse to QueueStore")
		}

		mi, ok := ps.Message.(*email.MessageInfo)
		if !ok {
			return nil, errors.New("Error in pull queue, bad parse to MessageInfo")
		}

		if err := mi.ParseBody(); err != nil {
			log.Printf("Body parrser error: %s", err.Error())
		}

		list = append(list, mi)
	}

	return list, nil
}

func (q *QueueBox) pushToQueue(key string, message *email.MessageInfo) {
	pq := q.queueFactory.GetOrCreate(key)
	qs := &QueueStore{
		Message:  message,
		Priority: 1,
		Key:      message.MessageId(),
	}

	pq.Push(qs)
}

func (q *QueueBox) queueId(key, kind string) (string, error) {
	val := fmt.Sprintf("%s@%s", key, kind)

	hash, err := config.ComputeHash(val)
	if err != nil {
		return "", errors.New("Bad queue id")
	}

	return hash, nil
}

func (q *QueueBox) acquireEmail() (*email.Email, error) {
	e := q.emailPool.Get().(*email.Email)
	defer q.emailPool.Put(e)

	if err := e.Init(); err != nil {
		return nil, err
	}

	return e, nil
}
