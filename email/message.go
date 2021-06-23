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
	SenderHeader    = "From:"
	RecipientHeader = "To:"
	SubjectHeader   = "Subject:"
)

var (
	ErrSubjectExist = errors.New("Subject of the message is already exists")
	ErrNoRecipient  = errors.New("No recipient")
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

func (m *Message) SetSender(name, address string) {
	a := mail.Address{
		Name:    name,
		Address: address,
	}
	m.header.Set(SenderHeader, a.String())
}

func (m *Message) AddRecipient(recipient string) {
	m.header.Add(RecipientHeader, recipient)
}

func (m *Message) SetSubject(subject string) {
	m.header.Set(SubjectHeader, subject)
}

func (m *Message) AddContent(content *Content) {
	m.contents = append(m.contents, content)
}

func (m *Message) AttachFile(file *File) {
	m.files = append(m.files, file)
}

func (m *Message) values(headerType string) []string {
	return m.header.Values(headerType)
}

func (m *Message) Recipients() string {
	r := m.values(RecipientHeader)
	return strings.Join(r, ",")
}

func (m *Message) Sender() string {
	return m.header.Get(SenderHeader)
}

func (m *Message) Subject() string {
	return m.header.Get(SubjectHeader)
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

//Returns all attached files
func (m *Message) FileCount() int {
	return len(m.files)
}

//Checking if any file attached to the message
func (m *Message) IsFile() bool {
	if m.FileCount() > 0 {
		return true
	}

	return false
}

func (m *Message) write() (*bytes.Buffer, error) {
	b := &bytes.Buffer{}

	writer := bufio.NewWriter(b)
	text := textproto.NewWriter(writer)

	boundary := fmt.Sprintf("%s", m.boundary())

	if err := m.writeHeader(text); err != nil {
		return nil, fmt.Errorf("Header writer error due to: %s", err)
	}

	if m.IsFile() {
		text.PrintfLine("Content-Type: multipart/mixed; boundary=%s", boundary)

		if err := m.mixedBody(boundary, b); err != nil {
			return nil, err
		}

	} else {
		text.PrintfLine("Content-Type: multipart/alternative; boundary=%s", boundary)

		if err := m.alternativeBody(boundary, b); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (m *Message) writeHeader(writer *textproto.Writer) error {
	writer.PrintfLine("From: %s", m.Sender())
	writer.PrintfLine("To: %s", m.Recipients())
	writer.PrintfLine("Subject: %s", m.Subject())
	writer.PrintfLine("MIME-Version: 1.0")

	if _, err := fmt.Fprint(writer.W, "\r\n"); err != nil {
		return err
	}

	return nil
}

func (m *Message) mixedBody(boundary string, writer io.Writer) error {
	mw := multipart.NewWriter(writer)
	mw.SetBoundary(boundary)

	h := make(textproto.MIMEHeader)

	alternativeBoundary := m.boundary()
	h.Set("Content-Type", fmt.Sprintf("mutlipart/alternative; boundary=%s", alternativeBoundary))

	w, err := mw.CreatePart(h)
	if err != nil {
		return err
	}

	if cmw, err := m.putContent(alternativeBoundary, w); err == nil {
		if err := cmw.Close(); err != nil {
			return err
		}
	} else {
		return err
	}

	if err := m.writeFile(mw); err != nil {
		return err
	}

	return mw.Close()
}

func (m *Message) alternativeBody(boundary string, writer io.Writer) error {
	mw, err := m.putContent(boundary, writer)
	if err != nil {
		return err
	}

	return mw.Close()
}

func (m *Message) putContent(boundary string, writer io.Writer) (*multipart.Writer, error) {
	mw := multipart.NewWriter(writer)

	if err := mw.SetBoundary(boundary); err != nil {
		return nil, err
	}

	for _, c := range m.contents {
		var mimetype string
		var b *bytes.Buffer

		h := make(textproto.MIMEHeader)

		if c.HtmlType {
			mimetype = "text/html"
			h.Set("Content-Transfer-Encoding", "quoted-printable")

			b = bytes.NewBuffer(c.Data)
		} else {
			mimetype = "text/plain"
			h.Set("Content-Transfer-Encoding", "base64")

			encode := m.encode(c.Data)
			b = bytes.NewBufferString(encode)
		}

		h.Set("Content-Type", fmt.Sprintf("%s; charset=UTF-8", mimetype))

		w, err := mw.CreatePart(h)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(w, b); err != nil {
			return nil, err
		}

	}

	return mw, nil
}

func (m *Message) writeFile(writer *multipart.Writer) error {
	for _, file := range m.files {

		ext := filepath.Ext(file.Name)
		mimetype := mime.TypeByExtension(ext)

		h := make(textproto.MIMEHeader)

		h.Set("Content-Type", mimetype)
		h.Set("Content-Transfer-Encoding", "base64")
		h.Set("Content-Disposition", fmt.Sprintf("attachment; filename==?UTF-8?B?%s?=\r\n", m.encode([]byte(file.Name))))

		encode := m.encode(file.Data)
		b := bytes.NewBufferString(encode)

		w, err := writer.CreatePart(h)
		if err != nil {
			return err
		}

		if _, err := io.Copy(w, b); err != nil {
			return err
		}

	}

	return nil
}
