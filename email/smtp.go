package email

import "net/smtp"

type SMTPServer struct {
}

func (s *SMTPServer) LoginAuth(config *Config) smtp.Auth {
	return LoginAuth(config.Username, config.Password)
}

func (s *SMTPServer) PlainAuth(config *Config) smtp.Auth {
	return smtp.PlainAuth("", config.Username, config.Password, config.SMTP.Hostname)
}

func (s *SMTPServer) CRAMMD5Auth(config *Config) smtp.Auth {
	return smtp.CRAMMD5Auth(config.Username, config.Password)
}
