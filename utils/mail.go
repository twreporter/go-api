package utils

import (
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"time"

	log "github.com/Sirupsen/logrus"
	"twreporter.org/go-api/models"
)

// EmailStrategy defines an interface to send emails
type EmailStrategy interface {
	Send(sender *EmailContext, to, subject, body string) error
}

// EmailContext sends emails by the provided email strategy
type EmailContext struct {
	conf  models.EmailSettings
	email EmailStrategy
}

// Send sends email using the given strategy
func (s *EmailContext) Send(to, subject, body string) error {
	return s.email.Send(s, to, subject, body)
}

// NewEmailSender ...
func NewEmailSender(conf models.EmailSettings, email EmailStrategy) *EmailContext {
	return &EmailContext{conf, email}
}

// NewSMTPEmailSender use smtp email sending strategy to send email
func NewSMTPEmailSender(conf models.EmailSettings) *EmailContext {
	return &EmailContext{conf, &SMTPEmailSender{smtp.SendMail}}
}

// SMTPEmailSender is an email sending method
type SMTPEmailSender struct {
	send func(string, smtp.Auth, string, []string, []byte) error
}

// Send sends email using the SMTP
func (s *SMTPEmailSender) Send(sender *EmailContext, to, subject, body string) error {
	emailSettings := sender.conf

	if len(emailSettings.SMTPServer) == 0 {
		log.Info("utils.mail.send: SMTPServer is not set")
		return nil
	}

	log.WithFields(log.Fields{
		"to":          to,
		"subject":     subject,
		"body":        body,
		"emailConfig": emailSettings,
	}).Debug("utils.mail.send")

	fromMail := mail.Address{Name: emailSettings.FeedbackName, Address: emailSettings.SMTPUsername}
	toMail := mail.Address{Name: "", Address: to}

	addr := emailSettings.SMTPServer + ":" + emailSettings.SMTPPort
	auth := LoginAuth(emailSettings.SMTPUsername, emailSettings.SMTPPassword)

	message := buildMessage(fromMail.String(), toMail.String(), subject, body)

	err := s.send(addr, auth, emailSettings.SMTPUsername, []string{to}, []byte(message))

	if err != nil {
		return models.NewAppError("Send", "utils.mail.send_mail", err.Error(), 500)
	}

	return nil
}

func encodeRFC2047Word(s string) string {
	return mime.BEncoding.Encode("utf-8", s)
}

func buildMessage(from, to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = encodeRFC2047Word(subject)
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	headers["Content-Transfer-Encoding"] = "8bit"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n<html><body>" + body + "</body></html>"
	return message
}

// loginAuth is used to implement smtp.Auth interface
type loginAuth struct {
	username, password string
}

// LoginAuth ...
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start begins an authentication with a server.
// It returns the name of the authentication protocol
// and optionally data to include in the initial AUTH message
// sent to the server. It can return proto == "" to indicate
// that the authentication should be skipped.
// If it returns a non-nil error, the SMTP client aborts
// the authentication attempt and closes the connection.
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next continues the authentication. The server has just sent
// the fromServer data. If more is true, the server expects a
// response, which Next should return as toServer; otherwise
// Next should return toServer == nil.
// If Next returns a non-nil error, the SMTP client aborts
// the authentication attempt and closes the connection.
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, models.NewAppError("Next", "utils.mail.smtp_authentication_machanism", "SMTP client aborts the authentication attempt and close the connection", 500)
		}
	}
	return nil, nil
}
