package email

import "errors"

type MessageInfo struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

type Stat struct {
	Key           string `json:"access_key"`
	MessageNumber int    `json:"message_number"`
	ID            int    `json:"id"`
}

func (m *MessageInfo) FindAll() ([]*MessageInfo, error) {
	return nil, errors.New("Not yet implemented")
}

func (m *MessageInfo) FindBySender(sender string) ([]*MessageInfo, error) {
	return nil, errors.New("Not yet implemented")

}

func (m *MessageInfo) FindBySubject(subject string) ([]*MessageInfo, error) {
	return nil, errors.New("Not yet implemented")
}
