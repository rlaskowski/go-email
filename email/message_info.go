package email

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/paulrosania/go-charset/charset"
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
	HtmlType bool   `json:"html"`
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
		reader:  reader,
		message: message,
	}

	return m
}

func (m *MessageInfo) Sender() *mail.Address {
	from := m.message.Header.Get("From")

	d, err := m.decode(from)
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

	dec, err := m.decode(s)
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

	m.contents = make([]*Content, 0)

	ct := part.Header.Get("Content-Type")

	reader, err := m.bodyReader(part, ct)
	if err != nil {
		return nil
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		ct = part.Header.Get("Content-Type")

		mediatype, _, err := mime.ParseMediaType(ct)
		if err != nil {
			return nil
		}

		c := &Content{}

		if strings.Contains(mediatype, "text/html") {
			c.HtmlType = true
		}

		dec, err := m.decodePart(part)
		if err != nil {
			return err
		}

		c.Data = dec

		m.contents = append(m.contents, c)

	}

	return nil
}

func (m *MessageInfo) writeFile(part *multipart.Part) error {
	if len(part.FileName()) > 0 {
		m.files = make([]*File, 0)

		filedata, err := m.decodePart(part)
		if err != nil {
			return err
		}

		filename, err := m.decode(part.FileName())
		if err != nil {
			filename = part.FileName()
		}

		file := &File{
			Name: filename,
			Data: filedata,
		}

		m.files = append(m.files, file)
	}

	return nil

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

func (m *MessageInfo) decode(encoded string) (string, error) {
	wd := new(mime.WordDecoder)

	wd.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if !strings.Contains(charset, "utf-8") {
			c, err := m.decodeCharset(charset, input)
			if err != nil {
				return nil, err
			}

			return bytes.NewReader(c), nil
		}

		return input, nil
	}

	return wd.DecodeHeader(encoded)

}

func (m *MessageInfo) decodeCharset(ch string, r io.Reader) ([]byte, error) {
	r, err := charset.NewReader(ch, r)
	if err != nil {
		return nil, err
	}

	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return c, nil
}
