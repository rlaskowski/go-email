package email

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/model"
	"github.com/rlaskowski/go-email/serialization"
)

type Email struct {
}

func NewEmail() *Email {
	return &Email{}
}

//Sending email
func (e *Email) Send(message *model.Message, file *model.File) error {
	c, err := e.loadConfig(message.Key)
	if err != nil {
		return fmt.Errorf("Could not load email config file due to: %s", err)
	}

	m := NewMessage()

	m.AddSender(message.Sender, c.Email)
	m.AddRecipient(message.Recipient)
	m.AddSubject(message.Subject)
	m.AddContent(message.Content)

	m.AttachFile(file)

	return e.send(c, m)
}

func (e *Email) AuthLoginAuth(config *Config) smtp.Auth {
	return AuthLoginAuth(config.Username, config.Password)
}

func (e *Email) PlainAuth(config *Config) smtp.Auth {
	return smtp.PlainAuth("", config.Username, config.Password, config.Hostname)
}

func (e *Email) CRAMMD5Auth(config *Config) smtp.Auth {
	return smtp.CRAMMD5Auth(config.Username, config.Password)
}

func (e *Email) send(config *Config, message *Message) error {
	auth := e.AuthLoginAuth(config)

	sender, err := mail.ParseAddress(message.Values(HeaderKeySender)[0])
	if err != nil {
		return err
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.Hostname, config.Port), auth, sender.Address, message.Values(HeaderKeyRecipient), message.Bytes())
	if err != nil {
		log.Printf("Error when try to send email due to: %s, client key %s", err, config.Key)
		return err
	}

	log.Printf("Email was sended successful to %s, client key %s", message.Recipients(), config.Key)

	return nil
}

func (e *Email) loadConfig(key string) (*Config, error) {
	var configList []*Config

	path := fmt.Sprintf("%s/config.yaml", config.GetWorkingDirectory())

	file, err := os.Open(path)
	defer file.Close()

	r, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := serialization.DeserializeFromYaml(&configList, r); err != nil {
		return nil, err
	}

	for _, c := range configList {
		if strings.Contains(key, c.Key) {
			return c, nil
		}
	}

	return nil, errors.New("Email configuration not found")

}
