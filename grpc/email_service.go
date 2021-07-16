package grpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

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

func (e *EmailQueue) MessageStat(request *emailservice.StatRequest, stream emailservice.EmailService_MessageStatServer) error {
	stat, err := e.emailServ.Stat(request.Key)
	if err != nil {
		return err
	}

	response := &emailservice.Stat{}

	for _, s := range stat {
		response.Key = s.Key
		response.MessageId = s.ID
		response.MessageNumber = s.MessageNumber

		stream.Send(response)
	}

	return nil
}

func (e *EmailQueue) ReceiveMessage(request *emailservice.IncomingMsgRequest, stream emailservice.EmailService_ReceiveMessageServer) error {
	m, err := e.emailServ.ReadMessage(request.Key, request.MessageNumber)
	if err != nil {
		return err
	}

	if err := m.ParseBody(); err != nil {
		return errors.New("Body parser error")
	}

	response := &emailservice.IncomingMsgResponse{}

	incomingMesssage := &emailservice.IncomingMessage{
		Id: m.MessageId(),
		Address: &emailservice.Address{
			Name:    m.Sender().Name,
			Address: m.Sender().Address,
		},
		Subject: m.Subject(),
		Date:    m.Date(),
	}

	for _, c := range m.Contents() {
		content := &emailservice.Content{
			HtmlType: c.HtmlType,
			Data:     c.Data,
		}
		incomingMesssage.Contents = append(incomingMesssage.Contents, content)
	}

	for _, f := range m.Files() {
		file := &emailservice.File{
			Name: f.Name,
			Data: f.Data,
		}
		incomingMesssage.Files = append(incomingMesssage.Files, file)
	}

	encode, err := json.Marshal(incomingMesssage)
	if err != nil {
		return err
	}

	b := make([]byte, 32*1024)
	buff := bytes.NewBuffer(encode)

	for {
		n, err := buff.Read(b)

		if n > 0 {
			response.Message = b[0:n]
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		stream.Send(response)
	}

	return nil
}

func (e *EmailQueue) RouteMessage(stream emailservice.EmailService_RouteMessageServer) error {
	/* statlist := make(map[int64]*email.Stat)

	for {
		statin, err := stream.Recv()

		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}




	} */
	return nil
}

func (e *EmailQueue) emailQueue(key string) queue.QueueProcess {
	que, err := e.queueFactory.GetOrCreate(key)
	if err != nil {
		return nil
	}

	return que
}
