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

func (e *EmailQueue) ReceiveMessage(request *emailservice.ReceiveMessageRequest, stream emailservice.EmailService_ReceiveMessageServer) error {

	if request.MessageNumber > 0 {
		stat := email.Stat{Key: request.Key, MessageNumber: int(request.MessageNumber)}

		if err := e.receiveMessage(stat, stream); err != nil {
			return err
		}

	} else {

		err := e.emailServ.ReceiveStat(func(stat email.Stat) error {

			if err := e.receiveMessage(stat, stream); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EmailQueue) receiveMessage(stat email.Stat, stream emailservice.EmailService_ReceiveMessageServer) error {
	resp := new(emailservice.ReceiveMessageResponse)

	info, err := e.emailServ.ReadMessage(stat.Key, stat.MessageNumber)

	if err != nil {
		return err
	}

	sender := info.Sender()

	message := emailservice.ReceiveMessage{
		Id: info.MessageId(),
		Address: &emailservice.Address{
			Name:    sender.Name,
			Address: sender.Address,
		},
		Subject: info.Subject(),
	}

	resp.ReceiveMessage = &message

	files, err := info.Files()
	if err != nil {
		return err
	}

	for _, f := range files {
		file := &emailservice.File{
			Name: f.Name,
			Data: f.Data,
		}

		resp.ReceiveMessage.Files = append(resp.ReceiveMessage.Files, file)

	}

	if err := stream.Send(resp); err != nil {
		return err
	}

	return nil
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
