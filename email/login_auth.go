package email

import (
	"errors"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) *loginAuth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unrecognized server in authlogin kind")
		}
	}
	return nil, nil
}
