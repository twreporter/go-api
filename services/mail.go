package services

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/mail"
	"net/smtp"
	"time"

	log "github.com/Sirupsen/logrus"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// MailService defines an interface to be implemented
type MailService interface {
	Send(to, subject, body string) error
}

// NewAmazonMailService returns a AamzonMailStrategy struct with required config
func NewAmazonMailService() MailService {
	return &AmazonMailStrategy{conf: globals.Conf.Email.Amazon}
}

// AmazonMailStrategy implements MailService interface
type AmazonMailStrategy struct {
	conf configs.AmazonConfig
}

// Send is a pointer receiver function of AmazonMailStrategy,
// which uses SES to send the mail
func (s *AmazonMailStrategy) Send(to, subject, body string) error {
	emailSettings := s.conf

	if len(emailSettings.Sender) == 0 {
		log.Info("utils.mail.send: Sender is not set")
		return nil
	}

	// Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(emailSettings.AwsRegion)},
	)
	if err != nil {
		return models.NewAppError("AmazonMailSender.Send", "cannot create a session to AWS", err.Error(), http.StatusInternalServerError)
	}

	// Create an SES client in the session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(emailSettings.Charset),
					Data:    aws.String(body), // HTML body
				},
				Text: &ses.Content{
					Charset: aws.String(emailSettings.Charset),
					Data:    aws.String(body), // text body, i.e., the email body for recipients with non-HTML email clients
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(emailSettings.Charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(emailSettings.Sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	log.WithFields(log.Fields{
		"to":          to,
		"subject":     subject,
		"body":        body,
		"emailConfig": emailSettings,
		"results":     result,
	}).Debug("utils.mail.send")

	// Display error messages if they occur.
	if err != nil {
		ec := ""
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				ec = ses.ErrCodeMessageRejected + ": "
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				ec = ses.ErrCodeMailFromDomainNotVerifiedException + ": "
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				ec = ses.ErrCodeConfigurationSetDoesNotExistException + ": "
			}
		}
		return models.NewAppError("AmazonMailSender.Send", "internal server error: fail to send email", ec+err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// NewSMTPMailService returns a SMTPMailStrategy struct with required config
func NewSMTPMailService() MailService {
	return &SMTPMailStrategy{conf: globals.Conf.Email.SMTP}
}

// SMTPMailStrategy implements EmailStrategy interface
type SMTPMailStrategy struct {
	conf configs.SMTPConfig
}

// Send is a pointer receiver function of SMTPMailStrategy,
// which uses smtp servers to send the mail
func (s *SMTPMailStrategy) Send(to, subject, body string) error {
	emailSettings := s.conf

	if len(emailSettings.Server) == 0 {
		log.Info("utils.mail.send: SMTPServer is not set")
		return nil
	}

	log.WithFields(log.Fields{
		"to":          to,
		"subject":     subject,
		"body":        body,
		"emailConfig": emailSettings,
	}).Debug("utils.mail.send")

	fromMail := mail.Address{Name: emailSettings.FeedbackName, Address: emailSettings.Username}
	toMail := mail.Address{Name: "", Address: to}

	addr := emailSettings.Server + ":" + emailSettings.Port
	auth := LoginAuth(emailSettings.Username, emailSettings.Password)

	message := buildMessage(fromMail.String(), toMail.String(), subject, body)

	err := smtp.SendMail(addr, auth, emailSettings.Username, []string{to}, []byte(message))

	if err != nil {
		return models.NewAppError("SMTPEmailSender.Send", "internal server error: fail to send email", err.Error(), http.StatusInternalServerError)
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
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}
