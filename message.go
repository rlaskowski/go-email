package email

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/mail"
	"net/textproto"
	"strings"
)

const (
	HeaderKeySender    = "From"
	HeaderKeyRecipient = "To"
	HeaderKeySubject   = "Subject"
	HeaderKeyContent   = "Content"
)

var (
	ErrSubjectExist = errors.New("Subject of the message is already exists")
	ErrNoRecipient  = errors.New("No recipient")
	ErrContentExist = errors.New("Content of the message is already exists")
)

type Message struct {
	header textproto.MIMEHeader
	buffer bytes.Buffer
}

func NewMessage() *Message {
	return &Message{
		header: make(textproto.MIMEHeader),
	}
}

func (m *Message) AddSender(name, address string) {
	a := mail.Address{
		Name:    name,
		Address: address,
	}
	m.header.Add(HeaderKeySender, a.String())
}

func (m *Message) AddRecipient(recipient string) {
	m.header.Add(HeaderKeyRecipient, recipient)
}

func (m *Message) AddSubject(subject string) {
	m.header.Add(HeaderKeySubject, subject)
}

func (m *Message) AddContent(content string) {
	m.header.Add(HeaderKeyContent, content)
}

func (m *Message) Values(headerType string) []string {
	s := m.header.Values(headerType)
	return s
}

func (m *Message) create(key string, writer *textproto.Writer) {
	s := m.header.Values(key)

	for k, v := range s {
		writer.PrintfLine("%s: %s", k, v)
	}
}

func (m *Message) setHeaderValue(header, value string) string {
	return fmt.Sprintf("%s: %s\r\n", header, value)
}

func (m *Message) Recipients() string {
	r := m.Values(HeaderKeyRecipient)
	return strings.Join(r, ",")
}

func (m *Message) Sender() string {
	s := m.Values(HeaderKeySender)
	if !(len(s) > 0) {
		return ""
	}
	return s[0]
}

func (m *Message) Subject() string {
	s := m.Values(HeaderKeySubject)
	if !(len(s) > 0) {
		return ""
	}
	return s[0]
}

func (m *Message) Content() string {
	c := m.Values(HeaderKeyContent)
	return strings.Join(c, " ")
}

func (m *Message) String() string {
	m.write()

	return m.buffer.String()
}

func (m *Message) Bytes() []byte {
	m.write()

	return m.buffer.Bytes()
}

func (m *Message) write() *textproto.Writer {
	bW := bufio.NewWriter(&m.buffer)
	writer := textproto.NewWriter(bW)

	writer.PrintfLine("Content-Type: text/html; charset='UTF-8'")
	writer.PrintfLine("From: %s", m.Sender())
	writer.PrintfLine("To: %s", m.Recipients())
	writer.PrintfLine("Subject: %s\r\n", m.Subject())
	writer.PrintfLine("%s", m.Content())

	return writer
}
