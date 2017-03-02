// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

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

func encodeRFC2047Word(s string) string {
	return mime.BEncoding.Encode("utf-8", s)
}

func connectToSMTPServer(config *models.Config) (net.Conn, *models.AppError) {
	var conn net.Conn
	var err error

	if config.EmailSettings.ConnectionSecurity == models.ConnSecurityTLS {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         config.EmailSettings.SMTPServer,
		}

		conn, err = tls.Dial("tcp", config.EmailSettings.SMTPServer+":"+config.EmailSettings.SMTPPort, tlsconfig)
		if err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.connect_smtp.open_tls.app_error", err.Error(), 500)
		}
	} else {
		conn, err = net.Dial("tcp", config.EmailSettings.SMTPServer+":"+config.EmailSettings.SMTPPort)
		if err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.connect_smtp.open.app_error", err.Error(), 500)
		}
	}
	return conn, nil
}

func newSMTPClient(conn net.Conn, config *models.Config) (*smtp.Client, *models.AppError) {
	c, err := smtp.NewClient(conn, config.EmailSettings.SMTPServer+":"+config.EmailSettings.SMTPPort)
	if err != nil {
		log.Error("utils.mail.new_client.open.error", err.Error())
		return nil, models.NewAppError("SendMail", "utils.mail.connect_smtp.open_tls.app_error", err.Error(), 500)
	}
	auth := smtp.PlainAuth("", config.EmailSettings.SMTPUsername, config.EmailSettings.SMTPPassword, config.EmailSettings.SMTPServer+":"+config.EmailSettings.SMTPPort)
	if config.EmailSettings.ConnectionSecurity == models.ConnSecurityTLS {
		if err = c.Auth(auth); err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.new_client.auth.app_error", err.Error(), 500)
		}
	} else if config.EmailSettings.ConnectionSecurity == models.ConnSecurityStarttls {
		log.Info("Send mail with STARTTLS")
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         config.EmailSettings.SMTPServer,
		}
		c.StartTLS(tlsconfig)
		if config.EmailSettings.SMTPServerOwner == models.Office360 {
			log.Info("SMTP server is owned by office360")
			log.Info("username:", config.EmailSettings.SMTPUsername)
			log.Info("password:", config.EmailSettings.SMTPPassword)
			err = c.Auth(LoginAuth(config.EmailSettings.SMTPUsername, config.EmailSettings.SMTPPassword))
		} else {
			log.Info("SMTP server is not owned by office360")
			err = c.Auth(auth)
		}
		if err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.new_client.auth.app_error", err.Error(), 500)
		}
	} else if config.EmailSettings.ConnectionSecurity == models.ConnSecurityPlain {
		// note: go library only supports PLAIN auth over non-tls connections
		if err = c.Auth(auth); err != nil {
			return nil, models.NewAppError("SendMail", "utils.mail.new_client.auth.app_error", err.Error(), 500)
		}
	}
	return c, nil
}

// SendMail send email to <to> address and the email subject/body will be subject and body
func SendMail(to, subject, body string) *models.AppError {
	return SendMailUsingConfig(to, subject, body, Cfg)
}

// SendMailUsingConfig will connect ot the SMTP server and create the connection to STMP server, and send the email to <to> address
func SendMailUsingConfig(to, subject, body string, config *models.Config) *models.AppError {

	log.Info("config:", config)

	if len(config.EmailSettings.SMTPServer) == 0 {
		return nil
	}

	log.WithFields(log.Fields{
		"to":      to,
		"subject": subject,
	}).Debug("utils.mail.send_mail.sending.debug")

	fromMail := mail.Address{Name: config.EmailSettings.FeedbackName, Address: config.EmailSettings.SMTPUsername}
	toMail := mail.Address{Name: "", Address: to}

	headers := make(map[string]string)
	headers["From"] = fromMail.String()
	headers["To"] = toMail.String()
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

	conn, err1 := connectToSMTPServer(config)
	if err1 != nil {
		return err1
	}
	defer conn.Close()

	c, err2 := newSMTPClient(conn, config)
	if err2 != nil {
		return err2
	}
	defer c.Quit()
	defer c.Close()

	if err := c.Mail(fromMail.Address); err != nil {
		return models.NewAppError("SendMail", "utils.mail.send_mail.from_address.app_error", err.Error(), 500)
	}

	if err := c.Rcpt(toMail.Address); err != nil {
		return models.NewAppError("SendMail", "utils.mail.send_mail.to_address.app_error", err.Error(), 500)
	}

	w, err := c.Data()
	if err != nil {
		return models.NewAppError("SendMail", "utils.mail.send_mail.msg_data.app_error", err.Error(), 500)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return models.NewAppError("SendMail", "utils.mail.send_mail.msg.app_error", err.Error(), 500)
	}

	err = w.Close()
	if err != nil {
		return models.NewAppError("SendMail", "utils.mail.send_mail.close.app_error", err.Error(), 500)
	}

	return nil
}
