package graphql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/twreporter/go-api/configs/constants"
	"github.com/twreporter/go-api/globals"
)

var client *Client
var sessionToken string

func NewClient() error {
	url := globals.Conf.MemberCMS.Url
	if len(url) == 0 {
		return errors.New("member cms url not set in config.go")
	}
	client = newClient(url)
	if globals.Conf.Environment == "development" {
		client.Log = func(s string) { log.Println(s) }
	}
	if err := refreshToken(); err != nil {
		return err
	}
	return nil
}

func Query(req *Request) (interface{}, error) {
	cookie := getCookie()
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Host", globals.Conf.MemberCMS.Host)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*constants.MemberCMSQueryTimeout)
	defer cancel()

	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	return respData, nil
}

func refreshToken() error {
	var respData interface{}

	req := NewRequest(`
    mutation Mutation($email: String!, $password: String!) {
  		authenticateSystemUserWithPassword(email: $email, password: $password) {
    		... on SystemUserAuthenticationWithPasswordSuccess {
      		sessionToken
    		}
    		... on SystemUserAuthenticationWithPasswordFailure {
      		message
    		}
  		}
		}
	`)
	req.Var("email", globals.Conf.MemberCMS.Email)
	req.Var("password", globals.Conf.MemberCMS.Password)
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("Host", globals.Conf.MemberCMS.Host)

	if err := client.Run(context.Background(), req, &respData); err != nil {
		return err
	}
	token, err := getValueFromField(respData, "sessionToken")
	if err != nil {
		return err
	}
	sessionToken = token
	return nil
}

func getValueFromField(source interface{}, field string) (string, error) {
	var value string
	var err error

	m, ok := source.(map[string]interface{})
	if !ok {
		return "", errors.New("type assertion failed")
	}
	for k, v := range m {
		if k == field {
			value = v.(string)
			break
		}
		value, err = getValueFromField(v, field)
	}
	return value, err
}

func getCookie() string {
	return fmt.Sprintf("keystonejs-session=%s", sessionToken)
}
