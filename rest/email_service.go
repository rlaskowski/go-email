package rest

import (
	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/queue"
)

type IncomingMessage struct {
	ID      string           `json:"ID"`
	Address Address          `json:"address"`
	Subject string           `json:"subject"`
	Date    string           `json:"date"`
	Content []*email.Content `json:"content"`
	File    []*email.File    `json:"file"`
}

type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type EmailService struct {
	queueBox *queue.QueueBox
}

func NewEmailService(queueBox *queue.QueueBox) *EmailService {
	return &EmailService{queueBox}
}

func (e *EmailService) ReceiveList(key string) ([]IncomingMessage, error) {
	qlist, err := e.queueBox.ReceiveMessage(key)
	if err != nil {
		return nil, err
	}

	var list []IncomingMessage

	for _, m := range qlist {
		im := IncomingMessage{
			ID:      m.MessageId(),
			Address: Address(*m.Sender()),
			Subject: m.Subject(),
			Date:    m.Date(),
		}

		copy(im.Content, m.Contents())
		copy(im.File, m.Files())

		list = append(list, im)
	}

	return list, nil
}
