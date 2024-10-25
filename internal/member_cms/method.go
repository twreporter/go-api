package member_cms

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/twreporter/go-api/configs/constants"
	"github.com/twreporter/go-api/globals"
)

type Session struct {
	token     string
	expiredAt time.Time
}

const graphqlEndpoint = "/graphql"

var client *Client
var session Session

func GetApiBaseUrl() (string, error) {
	url := globals.Conf.MemberCMS.Url
	if len(url) == 0 {
		return "", errors.New("member cms url not set in config.go")
	}
	return url, nil
}

func NewClient() error {
	if !globals.Conf.Features.MemberCMS {
		return errors.New("disable intergrating with member cms")
	}
	url, err := GetApiBaseUrl()
	if err != nil {
		return err
	}
	url = url + graphqlEndpoint

	client = newClient(url)
	if globals.Conf.Environment == "development" || globals.Conf.Environment == "staging" {
		client.Log = func(s string) { log.Println(s) }
	}

	// refresh token concurrently
	go refreshToken()

	return nil
}

func Query(req *Request) (interface{}, error) {
	var respData interface{}

	if !globals.Conf.Features.MemberCMS {
		return respData, errors.New("disable intergrating with member cms")
	}
	if err := getValidSession(); err != nil {
		return respData, err
	}

	cookie := getCookie()
	req.Header.Set("Cookie", cookie)
	req.Host = globals.Conf.MemberCMS.Host
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*constants.MemberCMSQueryTimeout)
	defer cancel()

	if err := client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	return respData, nil
}

func AppendRequiredHeader(req *http.Request) error {
	if err := getValidSession(); err != nil {
		return err
	}

	cookie := getCookie()
	req.Header.Set("Cookie", cookie)
	req.Host = globals.Conf.MemberCMS.Host

	return nil
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
	req.Host = globals.Conf.MemberCMS.Host

	if err := client.Run(context.Background(), req, &respData); err != nil {
		return err
	}
	token, err := getValueFromField(respData, "sessionToken")
	if err != nil {
		return err
	}
	session.token = token
	session.expiredAt = getExpiration()
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
	return fmt.Sprintf("keystonejs-session=%s", session.token)
}

// todo: get expiration from authenticate mutation response after member cms update auth api
func getExpiration() time.Time {
	return time.Now().Add(time.Second * time.Duration(globals.Conf.MemberCMS.SessionMaxAge))
}

func getValidSession() error {
	if session.token == "" || session.expiredAt.IsZero() || session.expiredAt.Before(time.Now()) {
		if err := refreshToken(); err != nil {
			return err
		}
	}

	return nil
}
