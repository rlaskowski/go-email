package pop3

import (
	"errors"
	"net"
	"net/mail"
	"net/textproto"
	"strconv"
	"strings"
)

const (
	OK    = "+OK"
	ERROR = "-ERR"
)

type Client struct {
	text       *textproto.Conn
	serverName string
}

func Dial(address string) (*Client, error) {
	conn, err := textproto.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(address)

	return NewClient(conn, host)
}

func NewClient(text *textproto.Conn, host string) (*Client, error) {
	_, err := readResponse(text)
	if err != nil {
		text.Close()
		return nil, err
	}

	c := &Client{
		text:       text,
		serverName: host,
	}

	return c, nil
}

func (c *Client) cmd(format string, args ...interface{}) (uint, error) {
	id, err := c.text.Cmd(format, args...)
	if err != nil {
		return 0, err
	}
	c.text.StartResponse(id)
	defer c.text.EndResponse(id)

	return id, nil
}

func (c *Client) Username(username string) error {
	if _, err := c.cmd("USER %s", username); err != nil {
		return err
	}

	if _, err := readResponse(c.text); err != nil {
		return err
	}

	return nil
}

func (c *Client) Password(password string) error {
	if _, err := c.cmd("PASS %s", password); err != nil {
		return err
	}

	if _, err := readResponse(c.text); err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(number int) error {
	if _, err := c.cmd("DELE %d", number); err != nil {
		return err
	}

	if _, err := readResponse(c.text); err != nil {
		return err
	}

	return nil
}

func (c *Client) Retr(number int) ([]string, error) {
	if _, err := c.cmd("RETR %d", number); err != nil {
		return nil, err
	}

	if _, err := readResponse(c.text); err != nil {
		return nil, mail.ErrHeaderNotPresent
	}

	return c.text.ReadDotLines()
}

func (c *Client) Stat() (string, error) {
	if _, err := c.cmd("STAT"); err != nil {
		return "", err
	}

	res, err := readResponse(c.text)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (c *Client) List() ([]string, error) {
	if _, err := c.cmd("LIST"); err != nil {
		return nil, err
	}

	if _, err := readResponse(c.text); err != nil {
		return nil, err
	}

	return c.text.ReadDotLines()
}

func (c *Client) Close() error {
	if _, err := c.cmd("QUIT"); err != nil {
		return err
	}

	if _, err := readResponse(c.text); err != nil {
		return err
	}

	return c.text.Close()
}

func (c *Client) Auth(username, password string) error {
	if err := c.Username(username); err != nil {
		return err
	}

	if err := c.Password(password); err != nil {
		return err
	}

	return nil
}

func (c *Client) MessageCount() int {
	stat, err := c.Stat()
	if err != nil {
		return 0
	}

	list := strings.Split(stat, " ")

	count, err := strconv.Atoi(list[0])
	if err != nil {
		return 0
	}

	return count
}

func (c *Client) readMessage() []*mail.Message {
	var messages []*mail.Message

	mcount := c.MessageCount()

	for m := 1; m == mcount; {

	}

	return messages

}

func readResponse(text *textproto.Conn) (string, error) {
	r, err := text.ReadLine()
	if err != nil {
		return "", err
	}

	res, err := parseResponse(r)
	if err != nil {
		return "", err
	}

	return res, nil

}

func parseResponse(line string) (string, error) {
	upperLine := strings.TrimSpace(strings.ToUpper(line))

	if strings.HasPrefix(upperLine, OK) {
		return strings.TrimSpace(line[3:]), nil
	}

	if strings.HasPrefix(upperLine, ERROR) {
		return "", errors.New(strings.TrimSpace(line[4:]))
	}

	return "", errors.New("Can not define response type")

}
