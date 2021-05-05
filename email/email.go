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
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/email/pop3"
	"github.com/rlaskowski/go-email/model"
	"github.com/rlaskowski/go-email/serialization"
)

type ReceiveFunc func(info Stat) error

type Email struct {
	config []*Config
	smtp   *SMTPServer
	mutex  *sync.Mutex
}

func NewEmail() *Email {
	return &Email{
		smtp:  new(SMTPServer),
		mutex: &sync.Mutex{},
	}
}

func (e *Email) Start() error {
	log.Print("Starting Email")
	return e.configure()
}

func (e *Email) Stop() error {
	log.Print("Stopping Email...")
	return nil
}

func (e *Email) configure() error {
	config, err := e.loadConfig()
	if err != nil {
		return err
	}

	e.config = config

	return nil
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

func (e *Email) ReceiveStat(fn ReceiveFunc) error {
	for _, c := range e.config {
		client, err := e.client(c.Key)
		defer client.Close()

		if err != nil {
			return err
		}

		if list, err := client.List(); err == nil {
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

	}

	return nil
}

func (e *Email) ReadMessage(key string, number int) (*MessageInfo, error) {
	client, err := e.client(key)
	defer client.Close()

	if err != nil {
		return &MessageInfo{}, err
	}

	r, err := client.Retr(number)
	if err != nil {
		return &MessageInfo{}, err
	}

	return NewMessageInfo(r), nil

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

func (e *Email) client(key string) (*pop3.Client, error) {
	c, err := e.configByKey(key)
	if err != nil {
		return &pop3.Client{}, err
	}

	address := net.JoinHostPort(c.POP3.Hostname, fmt.Sprintf("%d", c.POP3.Port))
	dial, err := pop3.Dial(address, c.POP3.Encryption)

	if err != nil {
		return &pop3.Client{}, err
	}

	if err := dial.Auth(c.Username, c.Password); err != nil {
		return &pop3.Client{}, err
	}

	return dial, nil

}

func (e *Email) loadConfig() ([]*Config, error) {
	var configList []*Config

	path := filepath.Join(config.GetWorkingDirectory(), "config.yaml")

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
