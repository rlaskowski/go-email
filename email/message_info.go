package email

import (
	"bytes"
	"encoding/base64"
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
	reader  textproto.Reader
	message *mail.Message
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

				filename, err := m.decode(part.FileName())
				if err != nil {
					filename = part.FileName()
				}

				files = append(files, File{Name: filename, Data: filedata})
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
