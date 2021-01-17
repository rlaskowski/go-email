package email

import (
	"errors"
	"net/textproto"
)

const (
	HeaderKeyFrom    = "from"
	HeaderKeyTo      = "to"
	HeaderKeySubject = "subject"
	HeaderKeyContent = "content"
)

var (
	ErrSubjectExist = errors.New("Subject of the message is already exists")
	ErrNoRecipient  = errors.New("No recipient")
	ErrContentExist = errors.New("Content of the message is already exists")
)

type Message struct {
	sender string
	header *textproto.MIMEHeader
}

func New(sender string) *Message {
	return &Message{
		sender: sender,
		header: new(textproto.MIMEHeader),
	}
}

func (m *Message) AddRecipient(recipient string) error {
	if !(len(recipient) > 0) {
		return ErrNoRecipient
	}

	m.header.Add(HeaderKeyTo, recipient)

	return nil
}

func (m *Message) AddSubject(subject string) error {
	s := m.header.Get(HeaderKeySubject)
	if len(s) > 0 {
		return ErrSubjectExist
	}

	m.header.Add(HeaderKeySubject, subject)

	return nil
}

func (m *Message) AddContent(content string) error {
	c := m.header.Get(HeaderKeyContent)
	if len(c) > 0 {
		return ErrContentExist
	}

	m.header.Add(HeaderKeyContent, content)

	return nil
}
