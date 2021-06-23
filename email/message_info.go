package email

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	_ "github.com/paulrosania/go-charset/data"
)

type MessageInfo struct {
	reader   *textproto.Reader
	message  *mail.Message
	files    []*File
	contents []*Content
}

type Stat struct {
	Key           string `json:"key"`
	MessageNumber int64  `json:"message_number"`
	ID            int64  `json:"message_id"`
}

type File struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type Content struct {
	HtmlType bool   `json:"html_type"`
	Data     []byte `json:"data"`
}

func NewMessageInfo(reader *textproto.Reader) *MessageInfo {
	return newMessageInfo(reader)
}

func newMessageInfo(reader *textproto.Reader) *MessageInfo {
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
		reader:   reader,
		message:  message,
		files:    make([]*File, 0),
		contents: make([]*Content, 0),
	}

	return m
}

func (m *MessageInfo) Sender() *mail.Address {
	from := m.message.Header.Get("From")

	d, err := decode(from)
	if err == nil {
		from = d
	}

	address, err := mail.ParseAddress(from)
	if err != nil {
		return &mail.Address{
			Name:    from,
			Address: from,
		}
	}

	return address
}

func (m *MessageInfo) Date() string {
	hdt, err := m.message.Header.Date()
	if err != nil {
		hdt = time.Now()
	}

	return hdt.Format(time.RFC3339Nano)
}

func (m *MessageInfo) Subject() string {
	s := m.message.Header.Get("Subject")

	dec, err := decode(s)
	if err != nil {
		return s
	}

	return dec
}

func (m *MessageInfo) MessageId() string {
	messageid := m.message.Header.Get("Message-ID")
	return strings.Trim(messageid, "<>")
}

func (m *MessageInfo) ParseBody() error {
	ct := m.message.Header.Get("Content-Type")

	reader, err := m.bodyReader(m.message.Body, ct)
	if err != nil {
		return err
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if err := m.writeContent(part); err != nil {
			return err
		}

		if err := m.writeFile(part); err != nil {
			return err
		}
	}

	return nil
}

func (m *MessageInfo) Contents() []*Content {
	return m.contents
}

func (m *MessageInfo) Files() []*File {
	return m.files
}

func (m *MessageInfo) writeContent(part *multipart.Part) error {
	if len(part.FileName()) > 0 {
		return nil
	}

	ct := part.Header.Get("Content-Type")

	reader, err := m.bodyReader(part, ct)
	if err != nil {
		return m.putContent(part)
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if err := m.putContent(part); err != nil {
			return nil
		}

	}

	return nil
}

func (m *MessageInfo) putContent(part *multipart.Part) error {
	ctp := part.Header.Get("Content-Type")

	mediatype, params, err := mime.ParseMediaType(ctp)
	if err != nil {
		return nil
	}

	charset := params["charset"]

	c := &Content{}

	dec, err := m.decodePart(part)
	if err != nil {
		return err
	}

	if strings.Contains(mediatype, "text/html") {
		c.HtmlType = true
	} else {

		r := bytes.NewReader(dec)

		dec, err = decodeCharset(strings.ToLower(charset), r)
		if err != nil {
			return err
		}

	}

	c.Data = dec

	m.contents = append(m.contents, c)

	return nil
}

func (m *MessageInfo) writeFile(part *multipart.Part) error {
	if !(len(part.FileName()) > 0) {
		return nil
	}

	filedata, err := m.decodePart(part)
	if err != nil {
		return err
	}

	filename, err := decode(part.FileName())
	if err != nil {
		filename = part.FileName()
	}

	file := &File{
		Name: filename,
		Data: filedata,
	}

	m.files = append(m.files, file)

	return nil

}

func (m *MessageInfo) isMultipart(mediatype string) bool {
	if !strings.HasPrefix(mediatype, "multipart") {
		return false
	}
	return true
}

func (m *MessageInfo) bodyReader(r io.Reader, v string) (*multipart.Reader, error) {
	mediatype, params, err := mime.ParseMediaType(v)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(mediatype, "multipart") {
		return nil, errors.New("No multipart in body")
	}

	reader := multipart.NewReader(r, params["boundary"])

	return reader, nil
}

func (m *MessageInfo) decodePart(part *multipart.Part) ([]byte, error) {
	var r io.Reader

	r = part

	cte := part.Header.Get("Content-Transfer-Encoding")

	if strings.Contains(cte, "base64") {
		r = base64.NewDecoder(base64.StdEncoding, part)
	}

	b := bytes.Buffer{}
	if _, err := b.ReadFrom(r); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
