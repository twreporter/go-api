package utils

import (
	"errors"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EmailTestSuite struct {
	suite.Suite
}

var settings SmtpEmailSettings

func (suite *EmailTestSuite) SetupTest() {
	settings = SmtpEmailSettings{
		SMTPUsername:       "fakeUser",
		SMTPPassword:       "fakePassword",
		SMTPServer:         "fakeServer",
		SMTPPort:           "fakePort",
		ConnectionSecurity: "STARTTLS",
		SMTPServerOwner:    "fakeOwner",
		FeedbackName:       "fakeName",
		FeedbackEmail:      "fakeAddress",
	}
}

func (suite *EmailTestSuite) TestSendMailFailure() {
	f := mockSend(errors.New("error"))
	ctx := NewEmailSender(&SMTPEmailSender{conf: settings, send: f})
	body := "Hello World"
	err := ctx.Send("receiver@twreporter.org", "mock subject", body)

	assert.NotNil(suite.T(), err)
}

func (suite *EmailTestSuite) TestSendMailSuccess() {
	f := mockSend(nil)
	ctx := NewEmailSender(&SMTPEmailSender{conf: settings, send: f})
	body := "Hello World"
	err := ctx.Send("receiver@twreporter.org", "mock subject", body)

	assert.Nil(suite.T(), err)
}

func mockSend(errToReturn error) func(string, smtp.Auth, string, []string, []byte) error {
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return errToReturn
	}
}

func TestEmailTestSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}
