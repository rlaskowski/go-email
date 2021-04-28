package grpc

import (
	"errors"

	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/grpc/protobuf/emailqueue"
	"github.com/rlaskowski/go-email/queue"
)

type EmailQueue struct {
	queueFactory *queue.QueueFactory
	emailServ    *email.Email
	emailqueue.UnimplementedEmailServiceServer
}

func NewEmailQueue(queueFactory *queue.QueueFactory, emserv *email.Email) *EmailQueue {
	return &EmailQueue{
		queueFactory: queueFactory,
		emailServ:    emserv,
	}
}

func (e *EmailQueue) ReceiveList(request *emailqueue.ReceiveRequest, stream emailqueue.EmailService_ReceiveListServer) error {
	resp := new(emailqueue.ReceiveStatResponse)

	err := e.receiveStatList(func(stat email.Stat) error {

		receiveStat := &emailqueue.ReceiveStat{Key: stat.Key, MessageNumber: int32(stat.MessageNumber), MessageId: int32(stat.ID)}
		resp.ReceiveStat = receiveStat

		if err := stream.Send(resp); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EmailQueue) ReceiveMessage(request *emailqueue.ReceiveRequest, stream emailqueue.EmailService_ReceiveMessageServer) error {
	resp := new(emailqueue.ReceiveMessageResponse)

	err := e.receiveStatList(func(stat email.Stat) error {

		info, err := e.emailServ.ReadMessage(request.Key, stat.MessageNumber)

		if err != nil {
			return err
		}

		sender := info.Sender()

		message := emailqueue.ReceiveMessage{
			Address: &emailqueue.Address{
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
			file := &emailqueue.File{
				Name: f.Name,
				Data: f.Data,
			}

			resp.ReceiveMessage.Files = append(resp.ReceiveMessage.Files, file)

		}

		if err := stream.Send(resp); err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

func (e *EmailQueue) receiveStatList(fn email.ReceiveFunc) error {
	que := e.emailQueue()

	sub, err := que.Subscribe(queue.SubjectReceiving)
	if err != nil {
		return nil
	}

	for _, s := range sub {
		stat, ok := s.Message[0].(email.Stat)
		if !ok {
			return errors.New("Could not convert Message to Stat")
		}

		fn(email.Stat{Key: stat.Key, MessageNumber: stat.MessageNumber, ID: stat.ID})
	}

	return nil

}

func (e *EmailQueue) emailQueue() queue.QueueConnection {
	que, err := e.queueFactory.GetOrCreate(queue.EmailQueueType)
	if err != nil {
		return nil
	}

	return que
}
