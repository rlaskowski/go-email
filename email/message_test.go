package email

import (
	"fmt"
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	senderName           = `Gopher`
	senderAddress        = `sender.gopher@golang.org`
	firstRecipientEmail  = `first.recipient@golang.org`
	secondRecipientEmail = `second.recipient@golang.org`
	subject              = `Golang message test`
)

func TestSender(t *testing.T) {
	m := NewMessage()
	m.SetSender(senderName, senderAddress)

	s := m.Sender()
	a, err := mail.ParseAddress(s)
	if err != nil {
		t.Errorf("Parse email error: %s", err)
	}

	assert.Equal(t, a.Name, senderName, fmt.Sprintf("Different sender name got %s expect %s", a.Name, senderName))
	assert.Equal(t, a.Address, senderAddress, fmt.Sprintf("Different sender address got %s expect %s", a.Address, senderAddress))
}

func TestSubject(t *testing.T) {
	m := NewMessage()
	m.SetSubject(subject)

	s := m.Subject()

	assert.Equal(t, s, subject, fmt.Sprintf("Different subject name got %s expect %s", s, subject))
}

func TestRecipients(t *testing.T) {
	m := NewMessage()
	m.AddRecipient(firstRecipientEmail)
	m.AddRecipient(secondRecipientEmail)

	assert.Containsf(t, m.Recipients(), firstRecipientEmail, "Recipient %s not found on recipients list", firstRecipientEmail)
	assert.Containsf(t, m.Recipients(), secondRecipientEmail, "Recipient %s not found on recipients list", secondRecipientEmail)

}

/* func TestWriter(t *testing.T) {
	data := make([]byte, 32*1024)

	rand.Read(data)
	c := &Content{
		HtmlType: false,
		Data:     data,
	}

    f := &File{
		Name: "example_file_name.pdf",
		Data: data,
	}

	m := NewMessage()

	m.AddSender("Sender Test", "test@test.com")
	m.AddRecipient("recipient@test.com")
	m.AddSubject("Test Subject")
	m.AddContent(c)

	//m.AttachFile(f)

	b, err := m.Bytes()
	if err != nil {
		t.Errorf("Failded read message due to: %s", err)
	}

} */
