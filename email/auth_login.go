package email

import (
	"errors"
	"net/smtp"
)

type authLogin struct {
	username, password string
}

func AuthLoginAuth(username, password string) *authLogin {
	return &authLogin{username, password}
}

func (a *authLogin) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *authLogin) Next(fromServer []byte, more bool) (toServer []byte, err error) {
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
