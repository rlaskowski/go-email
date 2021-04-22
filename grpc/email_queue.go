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
	return e.receiveStatList(stream)
}

func (e *EmailQueue) ReceiveMessage(request *emailqueue.ReceiveRequest, stream emailqueue.EmailService_ReceiveMessageServer) error {

	que := e.emailQueue()

	sub, err := que.Subscribe(queue.SubjectReceiving)
	if err != nil {
		return nil
	}

	resp := new(emailqueue.ReceiveMessageResponse)

	for _, s := range sub {
		stat, ok := s.Message[0].(email.Stat)
		if !ok {
			return errors.New("Could not convert Message to Stat")
		}

		info, err := e.emailServ.ReadMessage(request.Key, stat.MessageNumber)
		if err != nil {
			return err
		}

		message := emailqueue.ReceiveMessage{Sender: info.Sender(), Subject: info.Subject()}

		resp.ReceiveMessage = &message

		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

func (e *EmailQueue) receiveStatList(stream emailqueue.EmailService_ReceiveListServer) error {
	que := e.emailQueue()

	sub, err := que.Subscribe(queue.SubjectReceiving)
	if err != nil {
		return nil
	}

	resp := new(emailqueue.ReceiveStatResponse)

	for _, s := range sub {
		stat, ok := s.Message[0].(email.Stat)
		if !ok {
			return errors.New("Could not convert Message to Stat")
		}

		receiveStat := &emailqueue.ReceiveStat{Key: stat.Key, MessageNumber: int32(stat.MessageNumber), MessageId: int32(stat.ID)}
		resp.ReceiveStat = receiveStat

		if err := stream.Send(resp); err != nil {
			return err
		}
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
