package email

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
)

type MessageInfo struct {
	reader  textproto.Reader
	message *mail.Message
}

type Stat struct {
	Key           string `json:"access_key"`
	MessageNumber int64  `json:"message_number"`
	ID            int64  `json:"message_id"`
}

type File struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func NewMessageInfo(reader textproto.Reader) *MessageInfo {
	return newMessageInfo(reader)
}

func newMessageInfo(reader textproto.Reader) *MessageInfo {
	line, err := reader.ReadDotBytes()
	if err != nil {
		return nil
	}

	b := bytes.NewReader(line)

	message, err := mail.ReadMessage(b)
	if err != nil {
		return nil
	}

	m := &MessageInfo{
		reader:  reader,
		message: message,
	}

	return m
}

func (m *MessageInfo) Sender() *mail.Address {
	from := m.message.Header.Get("From")

	address, err := mail.ParseAddress(from)
	if err != nil {
		return nil
	}

	return address
}

func (m *MessageInfo) Subject() string {
	return m.message.Header.Get("Subject")
}

func (m *MessageInfo) MessageId() string {
	messageid := m.message.Header.Get("Message-ID")
	return strings.Trim(messageid, "<>")
}

func (m *MessageInfo) Files() ([]File, error) {
	files := make([]File, 0)

	ct := m.message.Header.Get("Content-Type")

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mediatype, "multipart") {
		reader := multipart.NewReader(m.message.Body, params["boundary"])

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}

			if len(part.FileName()) > 0 {
				filedata, err := m.decodeFile(part)
				if err != nil {
					return nil, err
				}

				files = append(files, File{Name: part.FileName(), Data: filedata})
			}

		}
	}

	return files, nil
}

func (m *MessageInfo) decodeFile(part *multipart.Part) ([]byte, error) {
	b := bytes.Buffer{}
	d := base64.NewDecoder(base64.StdEncoding, part)

	if _, err := b.ReadFrom(d); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
