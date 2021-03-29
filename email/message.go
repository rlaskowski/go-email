package email

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"net/textproto"
	"path/filepath"
	"strings"

	"github.com/rlaskowski/go-email/model"
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
	files  []*model.File
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

func (m *Message) AttachFile(file *model.File) {
	if file != nil {
		m.files = append(m.files, file)
	}
}

func (m *Message) Values(headerType string) []string {
	s := m.header.Values(headerType)
	return s
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

func (m *Message) encodeString(name string) string {
	encode := base64.StdEncoding.EncodeToString([]byte(name))
	return encode
}

func (m *Message) boundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", buf[:])
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

	boundary := fmt.Sprintf("%s", m.boundary())

	writer.PrintfLine("From: %s", m.Sender())
	writer.PrintfLine("To: %s", m.Recipients())
	writer.PrintfLine("Subject: %s", m.Subject())
	writer.PrintfLine("MIME-Version: 1.0")

	if len(m.files) > 0 {
		writer.PrintfLine("Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
		writer.PrintfLine("--%s", boundary)
	}

	writer.PrintfLine("Content-Type: text/plain; charset=utf-8\r\n")
	writer.PrintfLine("%s", m.Content())

	m.writeFile(boundary, writer)

	return writer
}

func (m *Message) writeFile(boundary string, writer *textproto.Writer) {
	if len(m.files) > 0 {
		for i, file := range m.files {
			writer.PrintfLine("\r\n\r\n--%s", boundary)

			ext := filepath.Ext(file.Name)
			mimetype := mime.TypeByExtension(ext)

			writer.PrintfLine("Content-Type: %s", mimetype)
			writer.PrintfLine("Content-Transfer-Encoding: base64")
			writer.PrintfLine("Content-Disposition: attachment; filename==?UTF-8?B?%s?=\r\n", m.encodeString(file.Name))

			writer.PrintfLine("%s", m.encodeString(string(file.Data)))

			if i == len(m.files) {
				writer.PrintfLine("\r\n--%s--", boundary)
			} else {
				writer.PrintfLine("\r\n--%s", boundary)
			}
		}
	}
}
