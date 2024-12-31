package utils

import (
	"net/smtp"
)

type EmailSender interface {
	SendSystemMessage(subject string, message string, email string, name string, attach string, fileNameAttach string) error
}

var sender EmailSender

func InitializeSender(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom string) {
	sender = &unsafeEmailSender{
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
		SmtpUser: smtpUser,
		SmtpPass: smtpPass,
		SmtpFrom: smtpFrom,
	}
}

func GetEmailSender() EmailSender {
	return sender
}

// Todo: document
type unsafeEmailSender struct {
	SmtpHost string
	SmtpPort string
	SmtpUser string
	SmtpPass string
	SmtpFrom string
}

func (s *unsafeEmailSender) SendSystemMessage(
	subject string,
	message string,
	email string,
	name string,
	attach string,
	fileNameAttach string,
) error {
	panic("method not available")
}

type UnencryptedAuth struct {
	smtp.Auth
}

func (a UnencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
