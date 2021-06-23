package email

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	senderName           = `Gopher`
	senderAddress        = `sender.gopher@golang.org`
	firstRecipientEmail  = `first.recipient@golang.org`
	secondRecipientEmail = `second.recipient@golang.org`
	contentText          = `Example text in content`
	fileName             = `Example-Filename.txt`
	fileText             = `Example text in file`
	subject              = `Golang message test`
)

func createTestMessage() (*Message, error) {
	m := NewMessage()

	m.SetSender(senderName, senderAddress)

	m.AddRecipient(firstRecipientEmail)
	m.AddRecipient(secondRecipientEmail)

	m.SetSubject(subject)

	cb := &bytes.Buffer{}

	cm := strings.NewReader(contentText)
	if _, err := io.Copy(cb, cm); err != nil {
		return nil, err
	}

	c := &Content{
		Data: cb.Bytes(),
	}

	m.AddContent(c)

	return m, nil
}

func TestSender(t *testing.T) {
	m := NewMessage()
	m.SetSender(senderName, senderAddress)

	s := m.Sender()
	a, err := mail.ParseAddress(s)
	if err != nil {
		t.Errorf("Parse email error: %s", err)
	}

	assert.Equal(t, a.Name, senderName, fmt.Sprintf("Different sender name got %s expected %s", a.Name, senderName))
	assert.Equal(t, a.Address, senderAddress, fmt.Sprintf("Different sender address got %s expected %s", a.Address, senderAddress))
}

func TestSubject(t *testing.T) {
	m := NewMessage()
	m.SetSubject(subject)

	s := m.Subject()

	assert.Equal(t, s, subject, fmt.Sprintf("Different subject name got %s expected %s", s, subject))
}

func TestRecipients(t *testing.T) {
	m := NewMessage()
	m.AddRecipient(firstRecipientEmail)
	m.AddRecipient(secondRecipientEmail)

	assert.Containsf(t, m.Recipients(), firstRecipientEmail, "Recipient %s not found on recipients list", firstRecipientEmail)
	assert.Containsf(t, m.Recipients(), secondRecipientEmail, "Recipient %s not found on recipients list", secondRecipientEmail)

}

func TestBoundary(t *testing.T) {
	m := NewMessage()
	assert.NotEmpty(t, m.boundary(), "Boundary value should not be empty")
}

func TestPutContent(t *testing.T) {
	m, err := createTestMessage()
	if err != nil {
		t.Errorf("Could not create test message due to: %s", err)
	}

	b := m.boundary()

	buff := &bytes.Buffer{}

	if mw, err := m.putContent(b, buff); err != nil {
		t.Errorf("Could not write content message due to: %s", err)
	} else {
		if err := mw.Close(); err != nil {
			t.Errorf("Could not close message: %s", err)
		}
	}

	mr := strings.NewReader(buff.String())
	mp := multipart.NewReader(mr, b)

	for {
		part, err := mp.NextPart()

		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		mf := &MessageInfo{}
		i, err := mf.decodePart(part)
		if err != nil {
			t.Error(err)
		}

		assert.Equalf(t, contentText, string(i), "Different content message got %s expected %s", string(i), contentText)

	}

}

func TestAlternativeBody(t *testing.T) {
	m, err := createTestMessage()
	if err != nil {
		t.Errorf("Could not create test message due to: %s", err)
	}

	b := m.boundary()

	buff := &bytes.Buffer{}

	if err := m.alternativeBody(b, buff); err != nil {
		t.Errorf("Bad alternative body due to: %s", err)
	}

	mr := multipart.NewReader(buff, b)

	for {
		_, err := mr.NextPart()

		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
	}

}

func TestMixedBody(t *testing.T) {
	m, err := createTestMessage()
	if err != nil {
		t.Errorf("Could not create test message due to: %s", err)
	}

	buff := &bytes.Buffer{}

	filedata := strings.NewReader(fileText)

	if _, err := io.Copy(buff, filedata); err != nil {
		t.Error(err)
	}

	f := &File{
		Name: fileName,
		Data: buff.Bytes(),
	}

	m.AttachFile(f)

	buff.Reset()

	b := m.boundary()

	if err := m.mixedBody(b, buff); err != nil {
		t.Errorf("Bad mixed body due to: %s", err)
	}

	mr := multipart.NewReader(buff, b)

	for {
		_, err := mr.NextPart()

		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

	}
}

func TestWriteMessage(t *testing.T) {
	m, err := createTestMessage()
	if err != nil {
		t.Errorf("Could not create test message due to: %s", err)
	}

	b, err := m.write()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(b.String())

	if _, err := mail.ReadMessage(b); err != nil {
		t.Errorf("Bad message format: %s", err)
	}
}
