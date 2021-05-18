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
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"path/filepath"
	"strings"
)

const (
	SenderHeader    = "From"
	RecipientHeader = "To"
	SubjectHeader   = "Subject"
	ContentHeader   = "Content"
)

var (
	ErrSubjectExist = errors.New("Subject of the message is already exists")
	ErrNoRecipient  = errors.New("No recipient")
	ErrContentExist = errors.New("Content of the message is already exists")
)

type Message struct {
	header   textproto.MIMEHeader
	files    []*File
	contents []*Content
}

func NewMessage() *Message {
	return &Message{
		header:   make(textproto.MIMEHeader),
		files:    make([]*File, 0),
		contents: make([]*Content, 0),
	}
}

func (m *Message) AddSender(name, address string) {
	a := mail.Address{
		Name:    name,
		Address: address,
	}
	m.header.Add(SenderHeader, a.String())
}

func (m *Message) AddRecipient(recipient string) {
	m.header.Add(RecipientHeader, recipient)
}

func (m *Message) AddSubject(subject string) {
	m.header.Add(SubjectHeader, subject)
}

func (m *Message) AddContent(content *Content) {
	m.contents = append(m.contents, content)
}

func (m *Message) AttachFile(file *File) {
	m.files = append(m.files, file)
}

func (m *Message) Values(headerType string) []string {
	return m.header.Values(headerType)
}

func (m *Message) setHeaderValue(header, value string) string {
	return fmt.Sprintf("%s: %s\r\n", header, value)
}

func (m *Message) Recipients() string {
	r := m.Values(RecipientHeader)
	return strings.Join(r, ",")
}

func (m *Message) Sender() string {
	s := m.Values(SenderHeader)
	return strings.Join(s, ",")
}

func (m *Message) Subject() string {
	s := m.Values(SubjectHeader)
	return strings.Join(s, " ")
}

func (m *Message) Content() string {
	c := m.Values(ContentHeader)
	return strings.Join(c, " ")
}

func (m *Message) encode(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func (m *Message) boundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", buf[:])
}

func (m *Message) String() (string, error) {
	b, err := m.write()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (m *Message) Bytes() ([]byte, error) {
	b, err := m.write()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (m *Message) write() (*bytes.Buffer, error) {
	b := &bytes.Buffer{}
	writer := textproto.NewWriter(bufio.NewWriter(b))

	writer.PrintfLine("From: %s", m.Sender())
	writer.PrintfLine("To: %s", m.Recipients())
	writer.PrintfLine("Subject: %s", m.Subject())
	writer.PrintfLine("MIME-Version: 1.0")

	if err := m.writeBody(writer); err != nil {
		return nil, err
	}

	/* 	if len(m.files) > 0 {
		writer.PrintfLine("Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
		writer.PrintfLine("--%s", boundary)
	}

	writer.PrintfLine("Content-Type: text/plain; charset=utf-8\r\n")
	writer.PrintfLine("%s", m.Content())

	if err := m.writeFile(boundary, writer); err != nil {
		return nil, err
	}*/

	return b, nil
}

func (m *Message) writeBody(writer *textproto.Writer) error {
	boundary := fmt.Sprintf("%s", m.boundary())

	/* mw := multipart.NewWriter(writer.W)

	if err := mw.SetBoundary(boundary); err != nil {
		return err
	}

	h := make(textproto.MIMEHeader) */
	multitype := "mixed"

	if len(m.files) > 0 {
		multitype = "alternative"
	}

	writer.PrintfLine("Content-Type: multipart/%s; boundary=%s", multitype, boundary)

	if err := m.writeContent(writer); err != nil {
		return err
	}

	return nil
}

func (m *Message) writeContent(writer *textproto.Writer) error {
	if len(m.files) > 0 {

	}

	return nil
}

func (m *Message) writeFile(boundary string, writer *textproto.Writer) error {
	if !(len(m.files) > 0) {
		return nil
	}

	mw := multipart.NewWriter(writer.W)

	if err := mw.SetBoundary(boundary); err != nil {
		return err
	}

	for _, file := range m.files {

		ext := filepath.Ext(file.Name)
		mimetype := mime.TypeByExtension(ext)

		h := make(textproto.MIMEHeader)

		h.Set("Content-Type", mimetype)
		h.Set("Content-Transfer-Encoding", "base64")
		h.Set("Content-Disposition", fmt.Sprintf("attachment; filename==?UTF-8?B?%s?=\r\n", m.encode([]byte(file.Name))))

		encode := m.encode(file.Data)
		b := bytes.NewBufferString(encode)

		w, err := mw.CreatePart(h)
		if err != nil {
			return err
		}

		if _, err := w.Write(b.Bytes()); err != nil {
			return err
		}

	}

	return mw.Close()

	/* writer.PrintfLine("\r\n\r\n--%s", boundary)



	writer.PrintfLine("Content-Type: %s", mimetype)
	writer.PrintfLine("Content-Transfer-Encoding: base64")
	writer.PrintfLine("Content-Disposition: attachment; filename==?UTF-8?B?%s?=\r\n", m.encodeString(file.Name))

	writer.PrintfLine("%s", m.encodeString(string(file.Data)))

	if i == len(m.files) {
		writer.PrintfLine("\r\n--%s--", boundary)
	} else {
		writer.PrintfLine("\r\n--%s", boundary)
	} */
}
