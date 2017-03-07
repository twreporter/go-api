package utils

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"time"
	"twreporter.org/go-api/models"
)

// EmailSender is an interface to define methods
type EmailSender interface {
	Send(to, subject, body string) *models.AppError
}

// NewSMTPEmailSender ...
func NewSMTPEmailSender(conf models.EmailSettings) EmailSender {
	return &smtpEmailSender{conf: conf, buildConn: ConnectToSMTPServer, newClient: NewSMTPClient, send: SendMailBySMTPClient}
}

type smtpEmailSender struct {
	conf      models.EmailSettings
	buildConn func(smtpServer, smtpPort, connSecurity string) (net.Conn, *models.AppError)
	newClient func(conn net.Conn, smtpServer, smtpPort, smtpUsername, smtpPassword, connSecurity, smtpServerOwner string) (*smtp.Client, *models.AppError)
	send      func(c *smtp.Client, from, to mail.Address, subject, body string) *models.AppError
}

// Implements EmailSender interface
func (sender *smtpEmailSender) Send(to, subject, body string) *models.AppError {
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

	conn, err := sender.buildConn(emailSettings.SMTPServer, emailSettings.SMTPPort, emailSettings.ConnectionSecurity)
	if err != nil {
		return err
	}

	defer conn.Close()

	c, err := sender.newClient(conn, emailSettings.SMTPServer, emailSettings.SMTPPort, emailSettings.SMTPUsername, emailSettings.SMTPPassword, emailSettings.ConnectionSecurity, emailSettings.SMTPServerOwner)

	if err != nil {
		return err
	}

	defer c.Quit()
	defer c.Close()

	fromMail := mail.Address{Name: emailSettings.FeedbackName, Address: emailSettings.SMTPUsername}
	toMail := mail.Address{Name: "", Address: to}

	err = sender.send(c, fromMail, toMail, subject, body)

	if err != nil {
		return err
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

// ConnectToSMTPServer create a connection to smtp server according to its connection security
func ConnectToSMTPServer(smtpServer, smtpPort, connectionSecurity string) (net.Conn, *models.AppError) {
	var conn net.Conn
	var err error

	if connectionSecurity == models.ConnSecurityTLS {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpServer,
		}

		conn, err = tls.Dial("tcp", smtpServer+":"+smtpPort, tlsconfig)
		if err != nil {
			return nil, models.NewAppError("connectToSMTPServer", "utils.mail.connect_smtp.open_tls.app_error", err.Error(), 500)
		}
	} else {
		conn, err = net.Dial("tcp", smtpServer+":"+smtpPort)
		if err != nil {
			return nil, models.NewAppError("connectToSMTPServer", "utils.mail.connect_smtp.open.app_error", err.Error(), 500)
		}
	}
	return conn, nil
}

// NewSMTPClient create a smtp client which has been authorized by smtp server
func NewSMTPClient(conn net.Conn, smtpServer, smtpPort, smtpUsername, smtpPassword, connectionSecurity, smtpServerOwner string) (*smtp.Client, *models.AppError) {
	c, err := smtp.NewClient(conn, smtpServer+":"+smtpPort)
	if err != nil {
		log.Error("utils.mail.new_client.open.error", err.Error())
		return nil, models.NewAppError("SendMail", "utils.mail.new_client.open.app_error", err.Error(), 500)
	}
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer+":"+smtpPort)
	if connectionSecurity == models.ConnSecurityTLS {
		if err = c.Auth(auth); err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.new_client.auth.app_error", err.Error(), 500)
		}
	} else if connectionSecurity == models.ConnSecurityStarttls {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpServer,
		}
		c.StartTLS(tlsconfig)
		if smtpServerOwner == models.Office360 {
			err = c.Auth(LoginAuth(smtpUsername, smtpPassword))
		} else {
			err = c.Auth(auth)
		}
		if err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.new_client.auth.app_error", err.Error(), 500)
		}
	} else if connectionSecurity == models.ConnSecurityPlain {
		// note: go library only supports PLAIN auth over non-tls connections
		if err = c.Auth(auth); err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.new_client.auth.app_error", err.Error(), 500)
		}
	}
	return c, nil
}

// SendMailBySMTPClient use smtp client to send email to the receiver
func SendMailBySMTPClient(c *smtp.Client, from, to mail.Address, subject, body string) *models.AppError {

	if c == nil {
		log.Info("utils.mail.send_mail: client is nil")
		return nil
	}

	message := buildMessage(from.String(), to.String(), subject, body)

	if err := c.Mail(from.Address); err != nil {
		return models.NewAppError("Send", "utils.mail.send.from_address.app_error", err.Error(), 500)
	}

	if err := c.Rcpt(to.Address); err != nil {
		return models.NewAppError("sendMailBy", "utils.mail.send.to_address.app_error", err.Error(), 500)
	}

	w, err := c.Data()
	if err != nil {
		return models.NewAppError("SendMail", "utils.mail.send.msg_data.app_error", err.Error(), 500)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return models.NewAppError("SendMail", "utils.mail.send.msg.app_error", err.Error(), 500)
	}

	err = w.Close()
	if err != nil {
		return models.NewAppError("SendMail", "utils.mail.send.close.app_error", err.Error(), 500)
	}

	return nil
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
