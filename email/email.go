package email

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

type Email struct {
	config Config
}

type Config struct {
	Hostname string
	Port     int
	Username string
	Password string
}

func NewEmail(config Config) *Email {
	return &Email{
		config: config,
	}
}

//Sending email
func (e *Email) Send(message *Message) error {
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Hostname)

	sender, err := mail.ParseAddress(message.Values(HeaderKeySender)[0])
	if err != nil {
		return err
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", e.config.Hostname, e.config.Port), auth, sender.Address, message.Values(HeaderKeyRecipient), message.Bytes())
	if err != nil {
		return err
	}
	return nil
}
