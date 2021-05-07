package grpc

import (
	"container/heap"
	"errors"

	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/grpc/protobuf/emailservice"
	"github.com/rlaskowski/go-email/queue"
)

type EmailQueue struct {
	queueFactory *queue.QueueFactory
	emailServ    *email.Email
	emailservice.UnimplementedEmailServiceServer
}

func NewEmailQueue(queueFactory *queue.QueueFactory, emserv *email.Email) *EmailQueue {
	return &EmailQueue{
		queueFactory: queueFactory,
		emailServ:    emserv,
	}
}

func (e *EmailQueue) ListMessages(request *emailservice.AuthRequest, stream emailservice.EmailService_ListMessagesServer) error {
	resp := new(emailservice.ReceiveMessage)

	err := e.emailServ.ReceiveStat(func(stat email.Stat) error {

		m, err := e.readMessage(stat)
		if err != nil {
			return err
		}

		resp = m

		stream.Send(resp)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EmailQueue) Message(request *emailservice.MessageRequest, stream emailservice.EmailService_MessageServer) error {
	resp := new(emailservice.ReceiveMessage)

	stat := email.Stat{Key: request.AuthRequest.Key, MessageNumber: request.MessageNumber}

	m, err := e.readMessage(stat)
	if err != nil {
		return err
	}

	resp = m

	stream.Send(resp)

	return nil
}

func (e *EmailQueue) readMessage(stat email.Stat) (*emailservice.ReceiveMessage, error) {
	info, err := e.emailServ.ReadMessage(stat.Key, stat.MessageNumber)

	if err != nil {
		return &emailservice.ReceiveMessage{}, err
	}

	sender := info.Sender()

	message := &emailservice.ReceiveMessage{
		Id: info.MessageId(),
		Address: &emailservice.Address{
			Name:    sender.Name,
			Address: sender.Address,
		},
		Subject: info.Subject(),
	}

	files, err := info.Files()
	if err != nil {
		return &emailservice.ReceiveMessage{}, err
	}

	for _, f := range files {
		file := &emailservice.File{
			Name: f.Name,
			Data: f.Data,
		}

		message.Files = append(message.Files, file)

	}

	return message, nil
}

func (e *EmailQueue) receiveStatList(fn email.ReceiveFunc) error {
	q := e.emailQueue()

	for q.Len() > 0 {
		qs := heap.Pop(q).(*queue.QueueStore)

		stat, ok := qs.Message.(email.Stat)
		if !ok {
			return errors.New("Could not convert Message to Stat")
		}

		fn(email.Stat{Key: stat.Key, MessageNumber: stat.MessageNumber, ID: stat.ID})

	}

	return nil

}

func (e *EmailQueue) emailQueue() queue.QueueProcess {
	que, err := e.queueFactory.GetOrCreate(queue.EmailQueueType)
	if err != nil {
		return nil
	}

	return que
}
