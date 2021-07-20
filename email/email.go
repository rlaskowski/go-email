package email

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/email/pop3"
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

func (e *Email) Init() error {
	config, err := e.loadConfig()
	if err != nil {
		return err
	}

	e.config = config

	return nil
}

//Sending email
func (e *Email) Send(key string, msg *Message) error {
	c, err := e.configByKey(key)
	if err != nil {
		return fmt.Errorf("Could not load email config file due to: %s", err)
	}

	sn := msg.SenderName()
	msg.SetSender(sn, c.Email)

	return e.send(c, msg)
}

func (e *Email) Stat(key string) ([]*Stat, error) {
	client, err := e.client(key)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	statList := make([]*Stat, 0)

	if list, err := client.List(); err == nil {

		for _, l := range list {
			m := strings.Split(l, " ")

			msgnumber, err := strconv.ParseInt(m[0], 0, 64)
			if err != nil {
				return nil, err
			}

			msgid, err := strconv.ParseInt(m[1], 0, 64)
			if err != nil {
				return nil, err
			}

			stat := &Stat{
				Key:           key,
				MessageNumber: msgnumber,
				ID:            msgid,
			}

			statList = append(statList, stat)
		}

	}

	return statList, nil

}

func (e *Email) ReadMessage(key string, number int64) (*MessageInfo, error) {
	client, err := e.client(key)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	r, err := client.Retr(int(number))
	if err != nil {
		return nil, err
	}

	return NewMessageInfo(r), nil
}

func (e *Email) send(config *Config, msg *Message) error {
	auth := e.smtp.LoginAuth(config)

	recipients := strings.Split(msg.Recipients(), ",")

	mb, err := msg.Bytes()
	if err != nil {
		return err
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.SMTP.Hostname, config.SMTP.Port), auth, msg.SenderAddress(), recipients, mb)
	if err != nil {
		log.Printf("Error when try to send email due to: %s, client key %s", err, config.Key)
		return err
	}

	log.Printf("Email was sended successful to %s, client key %s", msg.Recipients(), config.Key)

	return nil
}

func (e *Email) client(key string) (*pop3.Client, error) {
	c, err := e.configByKey(key)
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(c.POP3.Hostname, fmt.Sprintf("%d", c.POP3.Port))
	dial, err := pop3.Dial(address, c.POP3.Encryption)

	if err != nil {
		return nil, err
	}

	if err := dial.Auth(c.Username, c.Password); err != nil {
		return nil, err
	}

	return dial, nil

}

func (e *Email) loadConfig() ([]*Config, error) {
	var configList []*Config

	path := filepath.Join(config.GetWorkingDirectory(), config.EmailConfigFile)

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
