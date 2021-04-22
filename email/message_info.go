package email

import (
	"errors"
	"net/mail"
)

type MessageInfo struct {
	mail *mail.Message
}

type Stat struct {
	Key           string `json:"access_key"`
	MessageNumber int    `json:"message_number"`
	ID            int    `json:"id"`
}

func NewMessageInfo(mail *mail.Message) *MessageInfo {
	return &MessageInfo{mail}
}

func (m *MessageInfo) FindAll() ([]*MessageInfo, error) {
	return nil, errors.New("Not yet implemented")
}

func (m *MessageInfo) Sender() string {
	return m.mail.Header.Get("Sender")
}

func (m *MessageInfo) Subject() string {
	return m.mail.Header.Get("Subject")
}
