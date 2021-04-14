package email

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/email/pop3"
	"github.com/rlaskowski/go-email/model"
	"github.com/rlaskowski/go-email/serialization"
)

type ReceiveFunc func(info Stat) error

type Email struct {
	config []*Config
	smtp   *SMTPServer
	pop3   *pop3.Client
}

func NewEmail() *Email {
	c, err := loadConfig()
	if err != nil {
		c = make([]*Config, 0)
	}
	return &Email{
		config: c,
		smtp:   new(SMTPServer),
	}
}

//Sending email
func (e *Email) Send(message *model.Message, file *model.File) error {
	c, err := e.configByKey(message.Key)
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

func (e *Email) Receive(fn ReceiveFunc) error {
	for _, c := range e.config {
		address := net.JoinHostPort(c.POP3.Hostname, fmt.Sprintf("%d", c.POP3.Port))
		dial, err := pop3.Dial(address, c.POP3.Encryption)

		if err != nil {
			return fmt.Errorf("Could not connect to POP3 due to: %s", err.Error())
		}

		if err := dial.Auth(c.Username, c.Password); err != nil {
			return err
		}

		if list, err := dial.List(); err == nil {
			for _, l := range list {
				stat := strings.Split(l, " ")

				msgnumber, err := strconv.Atoi(stat[0])
				if err != nil {
					msgnumber = 0
				}

				msgid, err := strconv.Atoi(stat[1])
				if err != nil {
					msgid = 0
				}

				fn(Stat{Key: c.Key, MessageNumber: msgnumber, ID: msgid})
			}
		}

		dial.Close()

	}

	return nil
}

func (e *Email) ReadMessage(number int) (*mail.Message, error) {
	retr, err := e.pop3.Retr(number)
	if err != nil {
		return nil, err
	}

	m, err := mail.ReadMessage(retr.R)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (e *Email) receive(client *pop3.Client, fn ReceiveFunc) error {
	list, err := client.List()
	if err != nil {
		return err
	}
	fmt.Println(list)

	return nil
}

func (e *Email) send(config *Config, message *Message) error {
	auth := e.smtp.LoginAuth(config)

	sender, err := mail.ParseAddress(message.Values(HeaderKeySender)[0])
	if err != nil {
		return err
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.SMTP.Hostname, config.SMTP.Port), auth, sender.Address, message.Values(HeaderKeyRecipient), message.Bytes())
	if err != nil {
		log.Printf("Error when try to send email due to: %s, client key %s", err, config.Key)
		return err
	}

	log.Printf("Email was sended successful to %s, client key %s", message.Recipients(), config.Key)

	return nil
}

func loadConfig() ([]*Config, error) {
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

	return configList, nil

}

func (e *Email) configByKey(key string) (*Config, error) {
	for _, c := range e.config {
		if strings.Contains(key, c.Key) {
			return c, nil
		}
	}

	return nil, errors.New("Email configuration was not found")
}
